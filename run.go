package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/jkcfg/jk/pkg/deferred"
	"github.com/jkcfg/jk/pkg/resolve"
	"github.com/jkcfg/jk/pkg/std"

	v8 "github.com/ry/v8worker2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a jk program",
	Args:  runArgs,
	Run:   run,
}

type paramsOption struct {
	params *std.Params
	kind   std.ParamKind
}

func (p *paramsOption) String() string {
	return ""
}

func (p *paramsOption) Set(s string) error {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return errors.New("input parameters are of the form name=value")
	}
	path := parts[0]
	v := parts[1]

	switch p.kind {
	case std.ParamKindBoolean:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("could not parse '%s' as a boolean", v)
		}
		p.params.SetBool(path, b)
	case std.ParamKindNumber:
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("could not parse '%s' as a float64", v)
		}
		p.params.SetNumber(path, n)
	case std.ParamKindString:
		p.params.SetString(path, v)
	case std.ParamKindObject:
		o, err := std.NewParamsFromJSON(strings.NewReader(v))
		if err != nil {
			return fmt.Errorf("could not parse JSON '%s': %v", v, err)
		}
		p.params.SetObject(path, o)
	}

	return nil
}

func (p *paramsOption) Type() string {
	return "name=value"
}

var runOptions struct {
	outputDirectory string
	parameters      std.Params
}

func parameters(kind std.ParamKind) pflag.Value {
	return &paramsOption{
		params: &runOptions.parameters,
		kind:   kind,
	}
}

func init() {
	runOptions.parameters = std.NewParams()
	runCmd.PersistentFlags().StringVarP(&runOptions.outputDirectory, "output-directory", "o", "", "where to output generated files")
	runCmd.PersistentFlags().Var(parameters(std.ParamKindBoolean), "pb", "boolean input parameter")
	runCmd.PersistentFlags().Var(parameters(std.ParamKindNumber), "pn", "number input parameter")
	runCmd.PersistentFlags().Var(parameters(std.ParamKindString), "ps", "string input parameter")
	runCmd.PersistentFlags().Var(parameters(std.ParamKindObject), "po", "object input parameter")
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
		Parameters:      runOptions.parameters,
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
