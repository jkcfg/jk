package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jkcfg/jk/pkg/std"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Example: generateExamples(),
	Short:   "Generate configuration files",
	Args:    generateArgs,
	Run:     generate,
}

func generateExamples() string {
	b := bytes.Buffer{}
	b.WriteString("  specifying where are input files used by script and output generated files\n")
	b.WriteString("    jk generate -v -i ./inputdir -o ./outputdir ./scriptdir/script.js\n")
	b.WriteString("  specifying input parameters\n")
	b.WriteString("    jk generate -v -p path.k1.k2=value ./scriptdir/script.js\n")
	b.WriteString("  specifying input parameters and file containing parameters\n")
	b.WriteString("    jk generate -v -p key=value -f filename.json script.js\n")
	return b.String()
}

var generateOptions struct {
	vmOptions

	stdout bool
}

func init() {
	initAllVMFlags(generateCmd, &generateOptions.vmOptions)

	generateCmd.PersistentFlags().BoolVar(&generateOptions.stdout, "stdout", false, "print values on stdout")

	jk.AddCommand(generateCmd)
}

func skipException(err error) bool {
	return strings.Contains(err.Error(), "jk-internal-skip: ")
}

func generateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("generate requires an input script")
	}
	return nil
}

func generate(cmd *cobra.Command, args []string) {
	if generateOptions.inputDirectory == "" {
		filename := args[0]
		inputDir, err := filepath.Abs(filepath.Dir(filename))
		if err != nil {
			log.Fatal(err)
		}
		generateOptions.inputDirectory = inputDir
	}

	vm := newVM(&generateOptions.vmOptions, ".")
	vm.parameters.SetBool("jk.generate.stdout", generateOptions.stdout)

	if err := vm.Run("@jkcfg/std/cmd/<generate>", fmt.Sprintf(string(std.Module("cmd/generate-module.js")), args[0])); err != nil {
		if !skipException(err) {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}
