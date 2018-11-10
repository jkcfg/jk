package resolve

import (
	"fmt"
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
	loader    Loader
	base      string
	importers []Importer
}

// NewResolver creates a new Resolver.
func NewResolver(loader Loader, basePath string, importers ...Importer) *Resolver {
	return &Resolver{
		loader:    loader,
		base:      basePath,
		importers: importers,
	}
}

// ResolveModule imports the specifier from an import statement located in the
// referrer module.
func (r Resolver) ResolveModule(specifier, referrer string) int {
	// The first importer that resolver the specifier wins.
	var source string
	var candidates []string

	for _, importer := range r.importers {
		data, considered := importer.Import(r.base, specifier, referrer)
		candidates = append(candidates, considered...)
		if data != nil {
			source = string(data)
			break
		}
	}

	if source == "" {
		fmt.Fprintf(os.Stderr, "error: could not import '%s' from '%s'\n", specifier, filepath.Join(r.base, referrer))
		if len(candidates) > 0 {
			fmt.Fprintf(os.Stderr, "candidates considered:\n")
			for _, candidate := range candidates {
				fmt.Fprintf(os.Stderr, "    %s\n", candidate)
			}
		}
		return 1
	}

	resolver := r
	resolver.base = filepath.Dir(filepath.Join(r.base, specifier))
	if err := r.loader.LoadModule(specifier, source, resolver.ResolveModule); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return 1
	}
	return 0
}
