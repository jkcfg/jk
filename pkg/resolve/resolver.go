package resolve

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

// Resolver implements ES 2015 module resolving.
type Resolver struct {
	loader Loader
	base   string
}

// NewResolver creates a new Resolver.
func NewResolver(loader Loader, basePath string) *Resolver {
	return &Resolver{
		loader: loader,
		base:   basePath,
	}
}

// ResolveModule imports the specifier from an import statement located in the
// referrer module.
func (c Resolver) ResolveModule(specifier, referrer string) int {
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

	resolver := Resolver{loader: c.loader, base: filepath.Dir(path)}
	err = c.loader.LoadModule(specifier, string(codeBytes), resolver.ResolveModule)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return 1
	}
	return 0
}
