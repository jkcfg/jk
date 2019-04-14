package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/jkcfg/jk/pkg/deferred"
	"github.com/jkcfg/jk/pkg/resolve"
	"github.com/jkcfg/jk/pkg/std"

	v8 "github.com/jkcfg/v8worker2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a jk program",
	Args:  runArgs,
	Run:   run,
}

type paramSource int

const (
	paramSourceFile paramSource = iota
	paramSourceCommandLine
)

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
}

func (p *paramsOption) String() string {
	return ""
}

func (p *paramsOption) setFromFile(s string) error {
	params, err := std.NewParamsFromFile(s)
	if err != nil {
		return fmt.Errorf("%s: %v", s, err)
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
	verbose         bool
	outputDirectory string
	inputDirectory  string
	parameters      std.Params
}

func parameters(source paramSource) pflag.Value {
	return &paramsOption{
		params: &runOptions.parameters,
		source: source,
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
}

func (e *exec) onMessageReceived(msg []byte) []byte {
	return std.Execute(msg, e.worker, std.ExecuteOptions{
		Verbose:         runOptions.verbose,
		Parameters:      runOptions.parameters,
		OutputDirectory: runOptions.outputDirectory,
		Root:            std.ReadBase{Path: e.workingDir, Resources: e.resources},
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

	resolver := resolve.NewResolver(worker, scriptDir,
		&resolve.MagicImporter{Specifier: "@jkcfg/std/resource", Generate: resources.MakeModule},
		&resolve.StaticImporter{Specifier: "@jkcfg/std", Source: std.Module()},
		&resolve.FileImporter{},
		&resolve.NodeModulesImporter{ModuleBase: scriptDir},
	)
	if err := worker.LoadModule(path.Base(filename), string(input), resolver.ResolveModule); err != nil {
		log.Fatal(err)
	}
	deferred.Wait() // TODO(michael): hide this in std?
}
