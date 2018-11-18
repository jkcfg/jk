package main

import (
	"errors"
	"io/ioutil"
	"path"

	"github.com/dlespiau/jk/pkg/resolve"
	"github.com/dlespiau/jk/pkg/std"

	v8 "github.com/ry/v8worker2"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a jk program",
	Args:  runArgs,
	RunE:  run,
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

func run(cmd *cobra.Command, args []string) error {
	worker := v8.New(onMessageReceived)
	filename := args[0]
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	resolver := resolve.NewResolver(worker, ".",
		&resolve.StaticImporter{Specifier: "std", Source: std.Module()},
		&resolve.FileImporter{},
	)
	return worker.LoadModule(path.Base(filename), string(input), resolver.ResolveModule)
}
