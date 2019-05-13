package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/jkcfg/jk/pkg/deferred"
	"github.com/jkcfg/jk/pkg/record"
	"github.com/jkcfg/jk/pkg/resolve"
	"github.com/jkcfg/jk/pkg/std"
	v8 "github.com/jkcfg/v8worker2"
)

type vmOptions struct {
	verbose          bool
	outputDirectory  string
	inputDirectory   string
	parameters       std.Params
	parameterFiles   []string // list of files specified on the command line with -f.
	emitDependencies bool

	debugImports bool
}

func initVMFlags(cmd *cobra.Command, opts *vmOptions) {
	opts.parameters = std.NewParams()

	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().StringVarP(&opts.outputDirectory, "output-directory", "o", "", "where to output generated files")
	cmd.PersistentFlags().StringVarP(&opts.inputDirectory, "input-directory", "i", "", "where to find files read in the script; if not set, the directory containing the script is used")
	cmd.PersistentFlags().VarP(parameters(opts, paramSourceCommandLine), "parameter", "p", "set input parameters")
	parameterFlag := cmd.PersistentFlags().VarPF(parameters(opts, paramSourceFile), "parameters", "f", "load parameters from a JSON or YAML file")
	parameterFlag.Annotations = map[string][]string{
		cobra.BashCompFilenameExt: {"json", "yaml", "yml"},
	}
	cmd.PersistentFlags().BoolVarP(&opts.emitDependencies, "emit-dependencies", "d", false, "emit script dependencies")
	cmd.PersistentFlags().BoolVar(&opts.debugImports, "debug-imports", false, "trace import logic")
	cmd.PersistentFlags().MarkHidden("debug-imports")
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

type vm struct {
	vmOptions

	scriptDir string
	inputDir  string

	worker    *v8.Worker
	recorder  *record.Recorder
	resources *std.ModuleResources
}

func (vm *vm) onMessageReceived(msg []byte) []byte {
	return std.Execute(msg, vm.worker, std.ExecuteOptions{
		Verbose:         vm.verbose,
		Parameters:      vm.parameters,
		OutputDirectory: vm.outputDirectory,
		Root:            std.ReadBase{Path: vm.inputDir, Resources: vm.resources, Recorder: vm.recorder},
		DryRun:          vm.emitDependencies,
	})
}

func newVM(opts *vmOptions) *vm {
	vm := &vm{
		vmOptions: *opts,
		resources: std.NewModuleResources(),
	}

	if opts.emitDependencies {
		recorder := &record.Recorder{}
		// Add the parameter files to the list of dependencies.
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("run: unable to get current working directory:", err)
		}
		for _, f := range opts.parameterFiles {
			recorder.Record(record.ParameterFile, record.Params{
				"path": filepath.Join(cwd, f),
			})
		}
		vm.recorder = recorder
	}

	worker := v8.New(vm.onMessageReceived)
	if err := worker.Load("errorHandler", errorHandler); err != nil {
		log.Fatal(err)
	}
	if err := worker.Load("global", global); err != nil {
		log.Fatal(err)
	}
	vm.worker = worker

	resolve.Debug(opts.debugImports)

	return vm
}

func (vm *vm) SetWorkingDirectory(dir string) {
	vm.scriptDir = dir

	inputDir := dir
	if vm.inputDirectory != "" {
		var err error
		inputDir, err = filepath.Abs(vm.inputDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}
	vm.inputDir = inputDir
}

func (vm *vm) resolver() *resolve.Resolver {
	// TODO(damien): there's an ugly dependency here. The user of the vm object has
	// to call SetWorkingDir before being able to call Run* functions.
	resolver := resolve.NewResolver(vm.worker, vm.scriptDir,
		&resolve.MagicImporter{Specifier: "@jkcfg/std/resource", Generate: vm.resources.MakeModule},
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
		&resolve.NodeImporter{ModuleBase: vm.scriptDir},
	)
	resolver.SetRecorder(vm.recorder)
	return resolver
}

func (vm *vm) Run(specifier string, source string) error {
	resolver := vm.resolver()
	if err := vm.worker.LoadModule(specifier, source, resolver.ResolveModule); err != nil {
		return err
	}
	return vm.flush()
}

func (vm *vm) RunModule(specifier string, referrer string) error {
	resolver := vm.resolver()
	_, ret := resolver.ResolveModule(specifier, referrer)
	if ret != 0 {
		err := fmt.Errorf("unable to load module %q", specifier)
		return errors.Wrap(err, "run-module")
	}
	return vm.flush()
}

func (vm *vm) RunFile(filename string) error {
	// Add the script to the list of dependencies.
	if vm.recorder != nil {
		abspath, _ := filepath.Abs(filename)
		vm.recorder.Record(record.ImportFile, record.Params{
			"specifier": filename,
			"path":      abspath,
		})
	}
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	resolver := vm.resolver()
	if err := vm.worker.LoadModule(filepath.Base(filename), string(input), resolver.ResolveModule); err != nil {
		return err
	}

	return vm.flush()
}

func (vm *vm) flush() error {
	deferred.Wait() // TODO(michael): hide this in std?

	if vm.recorder != nil {
		data, err := json.MarshalIndent(vm.recorder, "", "  ")
		if err != nil {
			return errors.Wrap(err, "emit-dependencies")
		}
		fmt.Println(string(data))
	}

	return nil
}
