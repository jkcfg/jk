package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/ghodss/yaml"
	v8 "github.com/ry/v8worker2"
)

func goString(b []byte) string {
	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String()
}

// DELETE
func onMessageReceived(msg []byte) []byte {
	y, err := yaml.JSONToYAML([]byte(goString(msg)))
	if err != nil {
		log.Fatalf("yaml: %s", err)
		return nil
	}
	fmt.Print(string(y))
	return nil
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

func localLoadModule(worker *v8.Worker, specifier, referrer string, cb v8.ModuleResolverCallback) error {
	println("[DEBUG] load module specifier:", specifier, "; referrer:", referrer)

	path := specifier
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(referrer), specifier)
	}

	if filepath.Ext(path) == "" {
		_, err := os.Stat(path + ".js")
		switch {
		case os.IsNotExist(err):
			println("[DEBUG] no file at", path+".js")
			path = filepath.Join(path, "index.js")
		case err != nil:
			println("[ERROR] stat", path, ":", err.Error())
			return err
		default:
			path = path + ".js"
		}
	}

	println("[DEBUG] path:", path)

	// FIXME don't allow climbing out of the base directory with '../../...'
	if _, err := os.Stat(path); err != nil {
		println("[ERROR] error on stat", path, ":", err.Error())
		return err
	}
	codeBytes, err := ioutil.ReadFile(path)
	if err != nil {
		println("[ERROR] reading file", path, ":", err.Error())
		return err
	}
	err = worker.LoadModule(specifier, string(codeBytes), cb)
	if err != nil {
		println("[ERROR]", err.Error())
	}
	return err
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

	var resolve v8.ModuleResolverCallback
	resolve = func(specifier, referrer string) int {
		if err := localLoadModule(worker, specifier, referrer, resolve); err != nil {
			println("[ERROR]", err.Error())
			return 1
		}
		return 0
	}

	if err := worker.LoadModule(path.Base(filename), string(input), resolve); err != nil {
		log.Fatalf("error: %v", err)
	}
}
