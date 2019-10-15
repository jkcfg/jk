package main

import (
	"fmt"
	"strings"

	"github.com/jkcfg/jk/pkg/std"
	"github.com/spf13/pflag"
)

type paramSource int

const (
	paramSourceFile paramSource = iota
	paramSourceCommandLine
)

// paramsOption implements a pflag.Value for script input parameters.
type paramsOption struct {
	params *std.Params
	source paramSource
	files  *[]string
}

func parameters(opts *vmOptions, source paramSource) pflag.Value {
	return &paramsOption{
		params: &opts.parameters,
		source: source,
		files:  &opts.parameterFiles,
	}
}

func (p *paramsOption) String() string {
	return ""
}

func (p *paramsOption) setFromFile(s string) error {
	params, err := std.NewParamsFromFile(s)
	if err != nil {
		return fmt.Errorf("%s: %v", s, err)
	}
	if p.files != nil {
		*p.files = append(*p.files, s)
	}

	p.params.Merge(params)

	return nil
}

func (p *paramsOption) setFromCommandline(s string) error {
	parts := strings.SplitN(s, "=", 2)
	path := parts[0]
	v := parts[1]

	p.params.SetString(path, v)
	return nil
}

func (p *paramsOption) Set(s string) error {
	if p.source == paramSourceFile {
		return p.setFromFile(s)
	}
	return p.setFromCommandline(s)
}

func (p *paramsOption) Type() string {
	if p.source == paramSourceFile {
		return "filename"
	}
	return "name=value"
}
