package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jkcfg/jk/pkg/deferred"
	"github.com/jkcfg/jk/pkg/record"
	"github.com/jkcfg/jk/pkg/resolve"
	"github.com/jkcfg/jk/pkg/std"

	v8 "github.com/jkcfg/v8worker2"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Example: examples(),
	Short:   "Execute a jk program",
	Args:    runArgs,
	Run:     run,
}

// InlineSpecifier is used as the initial module specifier when exec'ing literal JavaScript
const InlineSpecifier = "<exec>"

// StdinSpecifier is used as the initial module specifier when reading JavaScript from stdin
const StdinSpecifier = "<stdin>"

// ToplevelReferrer is used as the module referrer when using --module
const ToplevelReferrer = "<toplevel>"

func examples() string {
	b := bytes.Buffer{}
	b.WriteString("  specifying where are input files used by script and output generated files\n")
	b.WriteString("    jk run -v -i ./inputdir -o ./outputdir ./scriptdir/script.js\n")
	b.WriteString("  specifying input parameters\n")
	b.WriteString("    jk run -v -p path.k1.k2=value ./scriptdir/script.js\n")
	b.WriteString("  specifying input parameters and file containing parameters\n")
	b.WriteString("    jk run -v -p key=value -f filename.json script.js\n")
	b.WriteString("  run the JavaScript given on the command line, with standard lib available\n")
	b.WriteString("    jk run -c log('foo')\n")
	b.WriteString("  run the module given, resolved relative to the current directory\n")
	b.WriteString("    jk run -m @example/module\n")
	b.WriteString("  read the script to run from stdin\n")
	b.WriteString("    jk run -\n")
	return b.String()
}

const errorHandler = `
function onerror(msg, src, line, col, err) {
  V8Worker2.print("Promise rejected at", src, line + ":" + col);
  V8Worker2.print(err.stack);
}
`

const inlineTemplate = `
import { log, write, read } from '@jkcfg/std';
import { dir, info } from '@jkcfg/std/fs';
import * as param from '@jkcfg/std/param';

%s;
`

const global = `
var global = {};
`

var runOptions struct {
	// control how the argument is interpreted; by default, it's a
	// file to load
	module, inline bool

	verbose          bool
	outputDirectory  string
	inputDirectory   string
	parameters       std.Params
	parameterFiles   []string // list of files specified on the command line with -f.
	emitDependencies bool

	debugImports bool
}

func init() {
	runOptions.parameters = std.NewParams()

	runCmd.PersistentFlags().BoolVarP(&runOptions.module, "module", "m", false, "treat argument as specifying a module to load")
	runCmd.PersistentFlags().BoolVarP(&runOptions.inline, "exec", "c", false, "treat argument as specifying literal JavaScript to execute")

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
	var (
		scriptDir string
	)

	// Before setting anything else up, we have to establish the
	// directory relative to which modules will be resolved.
	var err error
	switch {
	case runOptions.module && runOptions.inline:
		log.Fatal("supply one or neither of --module,-m and --exec,-c")
	case runOptions.module || runOptions.inline || args[0] == "-":
		scriptDir, err = filepath.Abs(".")
	default:
		filename := args[0]
		scriptDir, err = filepath.Abs(filepath.Dir(filename))
	}

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

	if runOptions.emitDependencies {
		engine.recorder = &record.Recorder{}
		// Add the parameter files to the list of dependencies.
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

	worker := v8.New(engine.onMessageReceived)
	engine.worker = worker

	if err := worker.Load("errorHandler", errorHandler); err != nil {
		log.Fatal(err)
	}
	if err := worker.Load("global", global); err != nil {
		log.Fatal(err)
	}

	resolve.Debug(runOptions.debugImports)
	resolver := resolve.NewResolver(worker, scriptDir,
		&resolve.MagicImporter{Specifier: "@jkcfg/std/resource", Generate: resources.MakeModule},
		&resolve.StaticImporter{Specifier: "std", Resolved: "@jkcfg/std/index.js", Source: std.Module("index.js")},
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
				ExternalName: "index.js", InternalModule: "index.js",
			}, {
				ExternalName: "param.js", InternalModule: "param.js",
			}, {
				ExternalName: "fs.js", InternalModule: "fs.js",
			}, {
				ExternalName: "merge.js", InternalModule: "merge.js",
			}},
		},
		&resolve.FileImporter{},
		&resolve.NodeImporter{ModuleBase: scriptDir},
	)
	resolver.SetRecorder(engine.recorder)

	var runErr error

	switch {
	case runOptions.module:
		_, ret := resolver.ResolveModule(args[0], ToplevelReferrer)
		if ret != 0 {
			runErr = fmt.Errorf("unable to load module %q", args[0])
		}
	case runOptions.inline:
		runErr = worker.LoadModule(InlineSpecifier, fmt.Sprintf(inlineTemplate, args[0]), resolver.ResolveModule)
	case args[0] == "-":
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		runErr = worker.LoadModule(StdinSpecifier, string(input), resolver.ResolveModule)
	default: // a file
		// Add the script to the list of dependencies.
		filename := args[0]
		if engine.recorder != nil {
			abspath, _ := filepath.Abs(filename)
			engine.recorder.Record(record.ImportFile, record.Params{
				"specifier": filename,
				"path":      abspath,
			})
		}
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		runErr = worker.LoadModule(filepath.Base(filename), string(input), resolver.ResolveModule)
	}

	if runErr != nil {
		log.Fatal(runErr)
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
