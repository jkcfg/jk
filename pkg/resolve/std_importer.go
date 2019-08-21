package resolve

import (
	"path/filepath"
	"strings"

	"github.com/jkcfg/jk/pkg/std"
)

const (
	stdPrefix = "@jkcfg/std"
)

// StdImporter is the standard library importer.
type StdImporter struct {
	// PublicModules are modules users are allowed to import.
	PublicModules []string
}

func isStdModule(name string) bool {
	return strings.HasPrefix(name, stdPrefix)
}

func (i *StdImporter) publicModule(path string) string {
	for _, m := range i.PublicModules {
		if path == m {
			return m
		}
	}
	return ""
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
		path = "index.js"
	}
	if !strings.HasSuffix(path, ".js") {
		path += ".js"
	}

	// fmt.Printf("import %s from %s (basePath=%s, path=%s)\n", specifier, referrer, basePath, path)

	// Ensure we only allow users to import PublicModules. Modules from the std lib
	// itself are allowed to import internal private modules.
	m := i.publicModule(path)
	if !isStdModule(referrer) && m == "" {
		trace(i, "'%s' is not a public module", specifier)
		return nil, "", candidate
	}

	source := path
	if isStdModule(basePath) {
		directory := basePath[len(stdPrefix):]
		source = filepath.Join(directory, path)
	}
	if m != "" {
		source = m
	}
	module := std.Module(source)
	if len(module) == 0 {
		trace(i, "'%s' is not part of the standard library", specifier)
		return nil, "", candidate
	}

	return module, filepath.Join(stdPrefix, source), candidate
}
