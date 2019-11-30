package resolve

import (
	"path"
	"strings"

	"github.com/jkcfg/jk/pkg/__std/lib"
	"github.com/jkcfg/jk/pkg/std"
	"github.com/jkcfg/jk/pkg/vfs"
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
func (i *StdImporter) Import(base vfs.Location, specifier, referrer string) ([]byte, vfs.Location, []Candidate) {
	candidate := []Candidate{{specifier, staticRule}}

	// Short circuit the lookup when:
	//  - we're not trying to load a @jkcfg/std.* module
	//  - we're not inside the std library resolution
	if !isStdModule(specifier) && !strings.HasPrefix(base.Path, stdPrefix) {
		return nil, vfs.Nowhere, candidate
	}

	p := specifier
	if isStdModule(p) {
		p = specifier[len(stdPrefix):]
	}
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		p = "index.js"
	}
	if !strings.HasSuffix(p, ".js") {
		p += ".js"
	}

	// fmt.Printf("import %s from %s (basePath=%s, path=%s)\n", specifier, referrer, basePath, path)

	// Ensure we only allow users to import PublicModules. Modules from the std lib
	// itself are allowed to import internal private modules.
	m := i.publicModule(p)
	if !isStdModule(referrer) && m == "" {
		trace(i, "'%s' is not a public module", specifier)
		return nil, vfs.Nowhere, candidate
	}

	sourcePath := p
	if isStdModule(base.Path) {
		directory := base.Path[len(stdPrefix):]
		sourcePath = path.Join(directory, p)
	}
	if m != "" {
		sourcePath = m
	}

	src := std.Module(sourcePath)
	if len(src) == 0 {
		trace(i, "'%s' is not part of the standard library", specifier)
		return nil, vfs.Nowhere, candidate
	}

	return src, vfs.Location{Vfs: lib.Assets, Path: path.Join(stdPrefix, sourcePath)}, candidate
}
