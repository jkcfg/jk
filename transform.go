package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/jkcfg/jk/pkg/std"
)

var transformCmd = &cobra.Command{
	Use:     "transform <script> <file>...",
	Example: transformExamples,
	Short:   "Transform configuration files",
	Args:    transformArgs,
	Run:     transform,
}

const transformExamples = `

  run the default export of script.js on each YAML value read from
  stdin, printing the transformed values to stdout

    jk transform ./script.js -

  run the default export of a module (or file) on each input document,
  and write the results back to outputdir/

    jk transform -o outputdir/ ./script.js ./inputdir/*.json

  run a function on each input file, and print the results to stdout

    jk transform --stdout -c '({ name: n, ...fields }) => ({ name: n + "-dev", ...fields })' inputdir/*.yaml
`

var transformOptions struct {
	vmOptions
	scriptOptions
	// print everything to stdout
	stdout bool
	// permit the overwriting of input files
	overwrite bool
	// force the format of the output
	format string
	// when reading from stdin (`-` is supplied as an argument),
	// expect a stream of values in the format given
	stdinFormat string
}

func init() {
	initScriptFlags(transformCmd, &transformOptions.scriptOptions)
	initExecFlags(transformCmd, &transformOptions.vmOptions)
	transformCmd.PersistentFlags().BoolVar(&transformOptions.stdout, "stdout", false, "print the resulting values to stdout (implied if reading from stdin)")
	transformCmd.PersistentFlags().BoolVar(&transformOptions.overwrite, "overwrite", false, "allow input file(s) to be overwritten by output file(s); otherwise, an error will be thrown")
	transformCmd.PersistentFlags().StringVar(&transformOptions.format, "format", "", "force all output values to this format")
	transformCmd.PersistentFlags().StringVar(&transformOptions.stdinFormat, "stdin-format", "yaml", "assume this format for values from stdin (either a 'yaml' stream or concatenated 'json' values)")
	jk.AddCommand(transformCmd)
}

func transformArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("transform requires an input script")
	}
	return nil
}

func transform(cmd *cobra.Command, args []string) {
	// We must use the current directory as the working directory (for
	// the purpose of resolving modules), because we're potentially
	// going to supply a path _relative to here_ as an import.
	vm := newVM(&transformOptions.vmOptions, ".")

	// Encode the inputs as a map of path to .. the same path (for
	// now). This is in part to get around the limitations of
	// parameters (arrays are not supported as values), and partly in
	// anticipation of there being more information to pass about each
	// input. Stdin is encoded as an empty string, since `"-"` could
	// be the name of a file.
	inputs := make(map[string]interface{})
	for _, f := range args[1:] {
		if f == "-" {
			inputs[""] = ""
			transformOptions.stdout = true
		} else {
			inputs[f] = f
		}
	}

	vm.parameters.Set("jk.transform.input", inputs)
	vm.parameters.Set("jk.transform.stdout", transformOptions.stdout)
	vm.parameters.Set("jk.transform.overwrite", transformOptions.overwrite)
	setGenerateFormat(transformOptions.format, vm)

	switch transformOptions.stdinFormat {
	case "json", "yaml":
		vm.parameters.Set("jk.transform.stdin.format", transformOptions.stdinFormat)
	default:
		log.Fatal("--stdin-format must be 'json' or 'yaml'")
	}

	var module string
	switch {
	case transformOptions.inline:
		module = fmt.Sprintf(string(std.Module("cmd/transform-exec.js")), args[0])
	default:
		module = fmt.Sprintf(string(std.Module("cmd/transform-module.js")), args[0])
	}
	if err := vm.Run("@jkcfg/std/cmd/<transform>", module); err != nil {
		log.Fatal(err)
	}
}
