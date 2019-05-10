package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

const inlineTemplate = `
import { log, write, read } from '@jkcfg/std';
import { dir, info } from '@jkcfg/std/fs';
import * as param from '@jkcfg/std/param';

%s;
`

var runOptions struct {
	vmOptions

	// control how the argument is interpreted; by default, it's a
	// file to load
	module, inline bool
}

func init() {
	runCmd.PersistentFlags().BoolVarP(&runOptions.module, "module", "m", false, "treat argument as specifying a module to load")
	runCmd.PersistentFlags().BoolVarP(&runOptions.inline, "exec", "c", false, "treat argument as specifying literal JavaScript to execute")

	initVMFlags(runCmd, &runOptions.vmOptions)

	jk.AddCommand(runCmd)
}

func runArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("run requires an input script")
	}
	return nil
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

	vm := newVM(&runOptions.vmOptions)
	vm.SetWorkingDirectory(scriptDir)

	var runErr error

	switch {
	case runOptions.module:
		runErr = vm.RunModule(args[0], ToplevelReferrer)
	case runOptions.inline:
		runErr = vm.Run(InlineSpecifier, fmt.Sprintf(inlineTemplate, args[0]))
	case args[0] == "-":
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		runErr = vm.Run(StdinSpecifier, string(input))
	default: // a file
		runErr = vm.RunFile(args[0])
	}

	if runErr != nil {
		log.Fatal(runErr)
	}

	vm.Wait()
	vm.Finish()
}
