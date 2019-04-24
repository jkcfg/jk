package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jkcfg/jk/pkg/deferred"
	"github.com/jkcfg/jk/pkg/record"
	"github.com/jkcfg/jk/pkg/resolve"
	"github.com/jkcfg/jk/pkg/std"

	v8 "github.com/jkcfg/v8worker2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Example: examples(),
	Short:   "Execute a jk program",
	Args:    runArgs,
	Run:     run,
}

type paramSource int

const (
	paramSourceFile paramSource = iota
	paramSourceCommandLine
)

func examples() string {
	b := bytes.Buffer{}
	b.WriteString("  specifying where are input files used by script and output generated files\n")
	b.WriteString("    jk run -v -i ./inputdir -o ./outputdir ./scriptdir/script.js\n")
	b.WriteString("  specifying input parameters\n")
	b.WriteString("    jk run -v -p path.k1.k2=value ./scriptdir/script.js\n")
	b.WriteString("  specifying input parameters and file containing parameters\n")
	b.WriteString("    jk run -v -p key=value -f filename.json script.js\n")
	return b.String()
}

const errorHandler = `
function onerror(msg, src, line, col, err) {
  V8Worker2.print("Promise rejected at", src, line + ":" + col);
  V8Worker2.print(err.stack);
}
`

const global = `
var global = {};
`

type paramsOption struct {
	params *std.Params
	source paramSource
	files  *[]string
}

func (p *paramsOption) String() string {
	return ""
}

func (p *paramsOption) setFromFile(s string) error {
	params, err := std.NewParamsFromFile(s)
	if err != nil {
		return fmt.Errorf("%s: %v", s, err)
	}
	if p.files != nil {
		*p.files = append(*p.files, s)
	}

	p.params.Merge(params)

	return nil
}

func (p *paramsOption) setFromCommandline(s string) error {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return errors.New("input parameters are of the form name=value")
	}
	path := parts[0]
	v := parts[1]

	p.params.SetString(path, v)
	return nil
}

func (p *paramsOption) Set(s string) error {
	if p.source == paramSourceFile {
		return p.setFromFile(s)
	}
	return p.setFromCommandline(s)
}

func (p *paramsOption) Type() string {
	if p.source == paramSourceFile {
		return "filename"
	}
	return "name=value"
}

var runOptions struct {
	verbose          bool
	outputDirectory  string
	inputDirectory   string
	parameters       std.Params
	parameterFiles   []string // list of files specified on the command line with -f.
	emitDependencies bool

	debugImports bool
}

func parameters(source paramSource) pflag.Value {
	return &paramsOption{
		params: &runOptions.parameters,
		source: source,
		files:  &runOptions.parameterFiles,
	}
}

func init() {
	runOptions.parameters = std.NewParams()
	runCmd.PersistentFlags().BoolVarP(&runOptions.verbose, "verbose", "v", false, "verbose output")
	runCmd.PersistentFlags().StringVarP(&runOptions.outputDirectory, "output-directory", "o", "", "where to output generated files")
	runCmd.PersistentFlags().StringVarP(&runOptions.inputDirectory, "input-directory", "i", "", "where to find files read in the script; if not set, the directory containing the script is used")
	runCmd.PersistentFlags().VarP(parameters(paramSourceCommandLine), "parameter", "p", "boolean input parameter")
	parameterFlag := runCmd.PersistentFlags().VarPF(parameters(paramSourceFile), "parameters", "f", "load parameters from a JSON or YAML file")
	parameterFlag.Annotations = map[string][]string{
		cobra.BashCompFilenameExt: {"json", "yaml", "yml"},
	}
	runCmd.PersistentFlags().BoolVarP(&runOptions.emitDependencies, "emit-dependencies", "d", false, "emit script dependencies")
	runCmd.PersistentFlags().BoolVar(&runOptions.debugImports, "debug-imports", false, "trace import logic")
	runCmd.PersistentFlags().MarkHidden("debug-imports")

	jk.AddCommand(runCmd)
}

func runArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("run requires an input script")
	}
	return nil
}

type exec struct {
	worker     *v8.Worker
	workingDir string
	resources  std.ResourceBaser
	recorder   *record.Recorder
}

func (e *exec) onMessageReceived(msg []byte) []byte {
	return std.Execute(msg, e.worker, std.ExecuteOptions{
		Verbose:         runOptions.verbose,
		Parameters:      runOptions.parameters,
		OutputDirectory: runOptions.outputDirectory,
		Root:            std.ReadBase{Path: e.workingDir, Resources: e.resources, Recorder: e.recorder},
		DryRun:          runOptions.emitDependencies,
	})
}

func run(cmd *cobra.Command, args []string) {
	filename := args[0]
	scriptDir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	inputDir := scriptDir
	if runOptions.inputDirectory != "" {
		inputDir, err = filepath.Abs(runOptions.inputDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}

	resources := std.NewModuleResources()

	engine := &exec{workingDir: inputDir, resources: resources}
	worker := v8.New(engine.onMessageReceived)
	engine.worker = worker

	if runOptions.emitDependencies {
		engine.recorder = &record.Recorder{}
	}

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := worker.Load("errorHandler", errorHandler); err != nil {
		log.Fatal(err)
	}
	if err := worker.Load("global", global); err != nil {
		log.Fatal(err)
	}

	resolve.Debug(runOptions.debugImports)
	resolver := resolve.NewResolver(worker, scriptDir,
		&resolve.MagicImporter{Specifier: "@jkcfg/std/resource", Generate: resources.MakeModule},
		&resolve.StaticImporter{Specifier: "std", Resolved: "@jkcfg/std/std.js", Source: std.Module("std.js")},
		&resolve.StdImporter{
			// List here the modules users are allowed to access. We map an external
			// module name to an internal module name to not link the file name used when
			// writing the standard library to a module name visible to the user.
			// eg.:
			//     import * as param from '@jkcfg/std/param.';
			// The name exposed to users is 'param.js', the file implementing this module
			// is 'std_param.js'
			//    { "param.js", "std_param.js" }
			PublicModules: []resolve.StdPublicModule{{
				ExternalName: "std.js", InternalModule: "std.js",
			}, {
				ExternalName: "param.js", InternalModule: "std_param.js",
			}, {
				ExternalName: "fs.js", InternalModule: "std_fs.js",
			}},
		},
		&resolve.FileImporter{},
		&resolve.NodeImporter{ModuleBase: scriptDir},
	)
	resolver.SetRecorder(engine.recorder)

	// Add the script and parameter files to the list of dependencies.
	if engine.recorder != nil {
		abspath, _ := filepath.Abs(filename)
		engine.recorder.Record(record.ImportFile, record.Params{
			"specifier": filename,
			"path":      abspath,
		})
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("run: unable to get current working directory:", err)
		}
		for _, f := range runOptions.parameterFiles {
			engine.recorder.Record(record.ParameterFile, record.Params{
				"path": filepath.Join(cwd, f),
			})
		}
	}

	if err := worker.LoadModule(path.Base(filename), string(input), resolver.ResolveModule); err != nil {
		log.Fatal(err)
	}
	deferred.Wait() // TODO(michael): hide this in std?

	if engine.recorder != nil {
		data, err := json.MarshalIndent(engine.recorder, "", "  ")
		if err != nil {
			log.Fatal("emit-dependencies:", err)
		}
		fmt.Println(string(data))
	}
}
