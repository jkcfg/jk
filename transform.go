package main

import (
	"fmt"

	"github.com/spf13/cobra"
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
}
