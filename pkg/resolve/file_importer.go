package resolve

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// FileImporter is an importer sourcing from a filesystem.
type FileImporter struct {
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
	if !filepath.IsAbs(path) {
		path = filepath.Join(basePath, specifier)
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

	// TODO don't allow climbing out of the base directory with '../../...'
	if _, err := os.Stat(path); err != nil {
		return nil, "", candidates
	}
	codeBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return codeBytes, path, candidates
}
