package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/dlespiau/jk/pkg/resolve"
	"github.com/dlespiau/jk/pkg/std"

	v8 "github.com/ry/v8worker2"
)

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: ")
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func onMessageReceived(msg []byte) []byte {
	return std.Execute(msg)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s INPUT", os.Args[0])
	}

	worker := v8.New(onMessageReceived)
	filename := os.Args[1]
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		errorf("%v", err)
	}

	resolver := resolve.NewResolver(worker, ".")
	if err := worker.LoadModule(path.Base(filename), string(input), resolver.ResolveModule); err != nil {
		errorf("%v", err)
	}
}
