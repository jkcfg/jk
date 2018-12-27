package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path"

	"github.com/jkcfg/jk/pkg/deferred"
	"github.com/jkcfg/jk/pkg/resolve"
	"github.com/jkcfg/jk/pkg/std"

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

type exec struct {
	worker *v8.Worker
}

func (e *exec) onMessageReceived(msg []byte) []byte {
	return std.Execute(msg, e.worker, std.ExecuteOptions{
		OutputDirectory: runOptions.outputDirectory,
	})
}

func run(cmd *cobra.Command, args []string) {
	engine := &exec{}
	worker := v8.New(engine.onMessageReceived)
	engine.worker = worker
	filename := args[0]
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	resolver := resolve.NewResolver(worker, path.Dir(filename),
		&resolve.StaticImporter{Specifier: "std", Source: std.Module()},
		&resolve.StaticImporter{Specifier: "@jkcfg/std", Source: std.Module()},
		&resolve.FileImporter{},
		&resolve.NodeModulesImporter{},
	)
	if err := worker.LoadModule(path.Base(filename), string(input), resolver.ResolveModule); err != nil {
		log.Fatal(err)
	}
	deferred.Wait() // TODO(michael): hide this in std?
}
