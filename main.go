package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/dlespiau/jk/pkg/std"

	v8 "github.com/ry/v8worker2"
)

func onMessageReceived(msg []byte) []byte {
	return std.Execute(msg)
}

// This is how the module loading works with V8Worker: You can ask a
// worker to load a module by calling `worker.LoadModule`. To this,
// you have to supply the specifier for the module (the name that was
// used to refer to it), the code in the module as a string, and a
// callback.
//
// The callback is used to load any modules imported in the code you
// provided; it's called with the specifier of the nested import, the
// referring module (our original specifier), and it's expected to
// load the imported module (i.e., by calling LoadModule itself). Some
// things are left implicit:
//
//  - there's no worker passed in the callback, so it has to be in the
//  closure, or otherwise accessed.
//
//  - the V8Worker code expects LoadModule to be called with the
//  specifier it gave, otherwise it will treat it as a failure to load
//  the module (NB this seems to mean you have to load a module
//  referred to by different paths once for each path)
//
//  - the referrer for an import will be the previous specifier; this
//  means you need to carry any directory context around with you,
//  since relative imports will otherwise lose the full path.

type resolveContext struct {
	worker *v8.Worker
	base   string
}

func (c resolveContext) resolveModule(specifier, referrer string) int {
	path := specifier
	if !filepath.IsAbs(path) {
		path = filepath.Join(c.base, specifier)
	}

	if filepath.Ext(path) == "" {
		_, err := os.Stat(path + ".js")
		switch {
		case os.IsNotExist(err):
			path = filepath.Join(path, "index.js")
		case err != nil:
			return 1
		default:
			path = path + ".js"
		}
	}

	// TODO don't allow climbing out of the base directory with '../../...'
	if _, err := os.Stat(path); err != nil {
		return 1
	}
	codeBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 1
	}

	resolver := resolveContext{worker: c.worker, base: filepath.Dir(path)}
	err = c.worker.LoadModule(specifier, string(codeBytes), resolver.resolveModule)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return 1
	}
	return 0
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s INPUT", os.Args[0])
	}

	worker := v8.New(onMessageReceived)
	filename := os.Args[1]
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	resolver := resolveContext{worker: worker, base: "."}
	if err := worker.LoadModule(path.Base(filename), string(input), resolver.resolveModule); err != nil {
		log.Fatalf("error: %v", err)
	}
}
