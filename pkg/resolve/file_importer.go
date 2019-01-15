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

// Import implements importer.
func (fi *FileImporter) Import(basePath, specifier, referrer string) ([]byte, string, []string) {
	var candidates []string

	path := specifier
	if !filepath.IsAbs(path) {
		path = filepath.Join(basePath, specifier)
	}

	if filepath.Ext(path) == "" {
		_, err := os.Stat(path + ".js")
		switch {
		case os.IsNotExist(err):
			candidates = append(candidates, path+".js")
			path = filepath.Join(path, "index.js")
		case err != nil:
			return nil, "", candidates
		default:
			path = path + ".js"
		}
	}

	candidates = append(candidates, path)

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
