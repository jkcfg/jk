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
`

var transformOptions struct {
	vmOptions
	scriptOptions
	stdout  bool // print everything to stdout
	inplace bool // permit the overwriting of input files
}

func init() {
	initScriptFlags(transformCmd, &transformOptions.scriptOptions)
	initVMFlags(transformCmd, &transformOptions.vmOptions)
	transformCmd.PersistentFlags().BoolVar(&transformOptions.stdout, "stdout", true, "print the resulting values to stdout")
	transformCmd.PersistentFlags().BoolVar(&transformOptions.inplace, "inplace", false, "allow input file(s) to be overwritten by output file(s)")
	jk.AddCommand(transformCmd)
}

func transformArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("transform requires an input script")
	}
	return nil
}

func transform(cmd *cobra.Command, args []string) {
	scriptDir := establishScriptDir(transformOptions.scriptOptions, args[0])
	vm := newVM(&transformOptions.vmOptions)
	vm.SetWorkingDirectory(scriptDir)

	// Encode the inputs as a map of path to .. the same path (for
	// now). This is in part to get around the limitations of
	// parameters (arrays are not supported as values), and partly in
	// anticipation of there being more information to pass about each
	// input.
	inputs := make(map[string]interface{})
	for _, f := range args[1:] {
		inputs[f] = f
	}
	vm.parameters.Set("jk.transform.input", inputs)

	if err := vm.Run("<transform>", fmt.Sprintf(string(std.Module("internal/transform.js")), args[0])); err != nil {
		log.Fatal(err)
	}
}
