package resolve

import (
	"path/filepath"
	"strings"

	"github.com/jkcfg/jk/pkg/std"
)

const (
	stdPrefix = "@jkcfg/std"
)

// StdPublicModule is a module exported by the standard library.
type StdPublicModule struct {
	// ExternalName is the name of the module as seen by the external world.
	ExternalName string
	// InternalModule is the file name of the module as embedded in the jk binary.
	InternalModule string
}

// StdImporter is the standard library importer.
type StdImporter struct {
	// PublicModules are modules users are allowed to import.
	PublicModules []StdPublicModule
}

func isStdModule(name string) bool {
	return strings.HasPrefix(name, stdPrefix)
}

func (i *StdImporter) publicModule(path string) *StdPublicModule {
	for _, m := range i.PublicModules {
		if path == m.ExternalName {
			return &m
		}
	}
	return nil
}

// Import implements importer.
func (i *StdImporter) Import(basePath, specifier, referrer string) ([]byte, string, []Candidate) {
	candidate := []Candidate{{specifier, staticRule}}

	// Short circuit the lookup when:
	//  - we're not trying to load a @jkcfg/std.* module
	//  - we're not inside the std library resolution
	if !isStdModule(specifier) && !strings.HasPrefix(basePath, stdPrefix) {
		return nil, "", candidate
	}

	path := specifier
	if isStdModule(path) {
		path = specifier[len(stdPrefix):]
	}
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		path = "std.js"
	}
	if !strings.HasSuffix(path, ".js") {
		path += ".js"
	}

	// fmt.Printf("import %s from %s (basePath=%s, path=%s)\n", specifier, referrer, basePath, path)

	// Ensure we only allow users to import PublicModules. Modules from the std lib
	// itself are allowed to import internal private modules.
	m := i.publicModule(path)
	if !isStdModule(referrer) && m == nil {
		trace(i, "'%s' is not a public module", specifier)
		return nil, "", candidate
	}

	source := path
	if m != nil {
		source = m.InternalModule
	}
	module := std.Module(source)
	if len(module) == 0 {
		trace(i, "'%s' is not part of the standard library", specifier)
		return nil, "", candidate
	}

	return module, filepath.Join(stdPrefix, source), candidate
}
