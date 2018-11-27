package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path"

	"github.com/justkidding-config/jk/pkg/resolve"
	"github.com/justkidding-config/jk/pkg/std"

	v8 "github.com/ry/v8worker2"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a jk program",
	Args:  runArgs,
	Run:   run,
}

var runOptions struct {
	outputDirectory string
}

func init() {
	runCmd.PersistentFlags().StringVarP(&runOptions.outputDirectory, "output-directory", "o", "", "where to output generated files")
	jk.AddCommand(runCmd)
}

func runArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("run requires an input script")
	}
	return nil
}

func onMessageReceived(msg []byte) []byte {
	return std.Execute(msg, std.ExecuteOptions{
		OutputDirectory: runOptions.outputDirectory,
	})
}

func run(cmd *cobra.Command, args []string) {
	worker := v8.New(onMessageReceived)
	filename := args[0]
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	resolver := resolve.NewResolver(worker, path.Dir(filename),
		&resolve.StaticImporter{Specifier: "std", Source: std.Module()},
		&resolve.FileImporter{},
	)
	if err := worker.LoadModule(path.Base(filename), string(input), resolver.ResolveModule); err != nil {
		log.Fatal(err)
	}
}
