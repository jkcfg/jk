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

// vmOptions are the options common (mostly) to all subcommands. Not
// all will be used (and thereby have a command-lin flag),
// necessarily; for example, `jk transform` doesn't use
// `inputDirectory`.
type vmOptions struct {
	verbose          bool
	outputDirectory  string
	inputDirectory   string
	parameters       std.Params
	parameterFiles   []string // list of files specified on the command line with -f.
	emitDependencies bool

	debugImports bool
}

// initInputFlags adds flags controlling input, to the given command
func initInputFlags(cmd *cobra.Command, opts *vmOptions) {
	cmd.PersistentFlags().StringVarP(&opts.inputDirectory, "input-directory", "i", "", "where to find files read in the script; if not set, the directory containing the script is used")
}

// initExecFlags adds flags controlling execution, to the given command
func initExecFlags(cmd *cobra.Command, opts *vmOptions) {
	opts.parameters = std.NewParams()

	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().StringVarP(&opts.outputDirectory, "output-directory", "o", "", "where to output generated files")
	cmd.PersistentFlags().VarP(parameters(opts, paramSourceCommandLine), "parameter", "p", "set input parameters")
	parameterFlag := cmd.PersistentFlags().VarPF(parameters(opts, paramSourceFile), "parameters", "f", "load parameters from a JSON or YAML file")
	parameterFlag.Annotations = map[string][]string{
		cobra.BashCompFilenameExt: {"json", "yaml", "yml"},
	}
	cmd.PersistentFlags().BoolVarP(&opts.emitDependencies, "emit-dependencies", "d", false, "emit script dependencies")
	cmd.PersistentFlags().BoolVar(&opts.debugImports, "debug-imports", false, "trace import logic")
	cmd.PersistentFlags().MarkHidden("debug-imports")
}

func initAllVMFlags(cmd *cobra.Command, opts *vmOptions) {
	initInputFlags(cmd, opts)
	initExecFlags(cmd, opts)
}

const errorHandler = `
function onerror(msg, src, line, col, err) {
  V8Worker2.log("Promise rejected at", src, line + ":" + col);
  V8Worker2.log(err.stack);
}
`

const global = `
var global = {};
`

func echo(args []interface{}) (interface{}, error) {
	// json.Marshal will serialise a []byte as base64-encoded;
	// stop it doing that by making all such args into []int
	// before responding.
	for i, arg := range args {
		if bytes, ok := arg.([]byte); ok {
			ints := make([]int, len(bytes), len(bytes))
			for j := range bytes {
				ints[j] = int(bytes[j])
			}
			args[i] = ints
		}
	}
	return args, nil
}

var rpcExtMethods = map[string]std.RPCFunc{
	"debug.echo": echo,
}

type vm struct {
	vmOptions

	scriptDir string
	inputDir  string

	worker    *v8.Worker
	recorder  *record.Recorder
	std       *std.Std
	resources *std.ModuleResources
}

func (vm *vm) onMessageReceived(msg []byte) []byte {
	return vm.std.Execute(msg, vm.worker)
}

func newVM(opts *vmOptions, workingDirectory string) *vm {
	vm := &vm{
		vmOptions: *opts,
		resources: std.NewModuleResources(),
	}

	/*
	 * Set scriptDir/inputDir based on workingDirectory and global options.
	 * This needs to be done early on as these values are used by both the
	 * stdlib and the module resolving mechanism.
	 */
	vm.setWorkingDirectory(workingDirectory)

	/* Setup a recorder object to gather the list of dependencies */
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

	/* create the stdlib */
	vm.std = std.NewStd(std.Options{
		Verbose:         vm.verbose,
		Parameters:      vm.parameters,
		OutputDirectory: vm.outputDirectory,
		Root:            std.ReadBase{Path: vm.inputDir, Resources: vm.resources, Recorder: vm.recorder},
		DryRun:          vm.emitDependencies,
		ExtMethods:      rpcExtMethods,
	})

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

func (vm *vm) setWorkingDirectory(dir string) {
	scriptDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	vm.scriptDir = scriptDir

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
		&resolve.StdImporter{
			// List here the modules users are allowed to access.
			PublicModules: []string{"index.js", "param.js", "fs.js", "merge.js", "debug.js", "schema.js"},
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
