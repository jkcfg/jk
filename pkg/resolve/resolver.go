package resolve

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/jkcfg/jk/pkg/record"
	"github.com/jkcfg/jk/pkg/vfs"
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

// Handy when debugging and implementing new importers
var debugImports bool

// Debug will cause the resolver code to trace what it's doing
func Debug(debug bool) {
	debugImports = debug
}

// ScriptBase returns a base location representing the filesystem
// under the path given. It is so called because it's used to create
// the initial base for resolving modules relative to the script being
// run.
func ScriptBase(path string) vfs.Location {
	return vfs.Location{Vfs: http.Dir(path), Path: "/"}
}

// Resolver implements module resolution by deferring to the set of
// importers that it's given.
type Resolver struct {
	recorder  *record.Recorder
	loader    Loader
	base      vfs.Location
	importers []Importer
}

// SetRecorder instructs Resolver to record actions in the specified recoder.
// Call with nil to disable.
func (r *Resolver) SetRecorder(recorder *record.Recorder) {
	r.recorder = recorder
}

// NewResolver creates a new Resolver.
func NewResolver(loader Loader, base vfs.Location, importers ...Importer) *Resolver {
	return &Resolver{
		loader:    loader,
		base:      base,
		importers: importers,
	}
}

func importerName(i Importer) string {
	return strings.TrimSuffix(reflect.ValueOf(i).Elem().Type().String()[8:], "Importer")
}

func trace(i Importer, f string, args ...interface{}) {
	if !debugImports {
		return
	}
	msg := fmt.Sprintf(f, args...)
	log.Printf("debug: % 6s: %s", importerName(i), msg)
}

func isInternalImporter(i Importer) bool {
	name := importerName(i)
	return name == "Std" || name == "Magic" || name == "Static"
}

// ResolveModule imports the specifier from an import statement located in the
// referrer module.
func (r Resolver) ResolveModule(specifier, referrer string) (string, int) {
	// The first importer that resolves the specifier wins.
	var resolved vfs.Location
	var source string
	var candidates []Candidate

	for _, importer := range r.importers {
		data, loc, considered := importer.Import(r.base, specifier, referrer)

		if len(data) == 0 {
			trace(importer, "✘ import %s from %s (base=%s)", specifier, referrer, r.base.Path) // TODO give a full account of the path
		} else {
			if r.recorder != nil && !isInternalImporter(importer) {
				r.recorder.Record(record.ImportFile, record.Params{
					"specifier": specifier,
					"path":      loc.Path,
				})
			}
			trace(importer, "✔ import %s from %s (base=%s) -> %s", specifier, referrer, r.base.Path, loc.Path) // TODO give a full account of the path
		}

		candidates = append(candidates, considered...)
		if data != nil {
			source = string(data)
			resolved = loc
			break
		}
	}

	if source == "" {
		fmt.Fprintf(os.Stderr, "error: could not import '%s' from '%s'\n", specifier, path.Join(r.base.Path, referrer))
		if len(candidates) > 0 {
			fmt.Fprintf(os.Stderr, "candidates considered:\n")
			for _, candidate := range candidates {
				fmt.Fprintf(os.Stderr, "    %s (%s)\n", candidate.Path, candidate.Rule)
			}
		}
		return "", 1
	}

	nextResolver := r
	nextResolver.base = vfs.Location{Vfs: resolved.Vfs, Path: path.Dir(resolved.Path)}
	// TODO the path will be used to uniquify modules, so it needs to
	// be uniquified itself, by the location, somehow
	if err := r.loader.LoadModule(resolved.Path, source, nextResolver.ResolveModule); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return "", 1
	}
	return resolved.Path, 0
}
