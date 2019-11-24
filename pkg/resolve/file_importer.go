package resolve

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileImporter is an importer sourcing from a filesystem, with the simple rules:
//  - an absolute path (starting with `/`) is not allowed
//  - a path starting with `./` or `../` is treated as a file relative
//  to the base path (previous resolution)
//  - any other path is treated as referring to a file relative to the
//  ModuleBase
//
//  A path `x/y/z` will resolve to
//   - the file `x/y/z`, if it exists
//   - the file `x/y/z.js` if it exists
//   - the file `x/y/z/index.js` if it exists
//  otherwise, nothing.
type FileImporter struct {
	ModuleBase string
}

const (
	expectedExtension   = ".js"
	extensionRulePrefix = "<path> -> <path>"
	extensionRule       = extensionRulePrefix + expectedExtension
	indexRulePrefix     = "<path> -> <path>/index"
	indexRule           = indexRulePrefix + expectedExtension
	verbatimRule        = "verbatim"
)

// Import implements importer.
func (fi *FileImporter) Import(basePath, specifier, referrer string) ([]byte, string, []Candidate) {
	var candidates []Candidate

	path := specifier
	if filepath.IsAbs(path) {
		log.Fatalf("absolute import path %q not allowed in %q", specifier, referrer)
	}

	// `import ... from 'foo/bar' -> treat as a module relative to ModuleBase
	// `import ... from './foo/bar' -> treat as a file relative to importer
	base := fi.ModuleBase
	if strings.HasPrefix(specifier, "./") || strings.HasPrefix(specifier, "../") {
		base = basePath
	}

	path = filepath.Join(base, specifier)
	rel, err := filepath.Rel(fi.ModuleBase, path)
	if err != nil {
		return nil, "", candidates
	}
	if strings.HasPrefix(rel, "../") { // outside the root
		return nil, "", candidates
	}

	if filepath.Ext(path) == "" {
		candidates = append(candidates, Candidate{path + ".js", extensionRule})
		_, err := os.Stat(path + ".js")
		switch {
		case os.IsNotExist(err):
			path = filepath.Join(path, "index.js")
			candidates = append(candidates, Candidate{path, indexRule})
		case err != nil:
			return nil, "", candidates
		default:
			path = path + ".js"
		}
	} else {
		candidates = append(candidates, Candidate{path, verbatimRule})
	}

	if _, err := os.Stat(path); err != nil {
		return nil, "", candidates
	}
	codeBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return codeBytes, path, candidates
}
