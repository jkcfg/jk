package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jkcfg/jk/pkg/std"
	"github.com/spf13/cobra"
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
}

func init() {
	initVMFlags(generateCmd, &generateOptions.vmOptions)

	jk.AddCommand(generateCmd)
}

func skipException(err error) bool {
	return strings.Contains(err.Error(), "throw new Error(\"jk-internal-skip: ")
}

func generateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("generate requires an input script")
	}
	return nil
}

func generate(cmd *cobra.Command, args []string) {
	filename := args[0]
	scriptDir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	vm := newVM(&generateOptions.vmOptions)
	vm.SetWorkingDirectory(scriptDir)

	if err = vm.Run("<generate>", fmt.Sprintf(string(std.Module("generate.js")), args[0])); err != nil {
		if !skipException(err) {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}
