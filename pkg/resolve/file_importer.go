package resolve

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/shurcooL/httpfs/vfsutil"

	"github.com/jkcfg/jk/pkg/vfs"
)

// FileImporter is an importer sourcing from a filesystem.  Modules
// are expected to be arranged at the root of the filesystem, in the
// directory structure implied by the import path. A specifier
// `foo/bar` will be resolved to (in the order attempted):
//
//   - `/foo/bar`
//   - `/foo/bar.js`
//   - `/foo/bar/index.js
type FileImporter struct {
	vfs http.FileSystem
}

// NewFileImporter constructs a FileImport given a filesystem
func NewFileImporter(vfs http.FileSystem) *FileImporter {
	return &FileImporter{vfs: vfs}
}

const (
	extensionRulePrefix = "<path> -> <path>"
	expectedExtension   = ".js"
	extensionRule       = extensionRulePrefix + expectedExtension
	indexRulePrefix     = "<path> -> <path>/index"
	indexRule           = indexRulePrefix + expectedExtension
	verbatimRule        = "verbatim"
)

// Import implements importer. Note that the file import only ever
// cares to look in the root of the given filesystem. It doesn't care
// where the importing module is located.
func (fi *FileImporter) Import(base vfs.Location, specifier, referrer string) ([]byte, vfs.Location, []Candidate) {
	if isRelative(specifier) {
		return nil, vfs.Nowhere, nil
	}
	return resolveFile(fi.vfs, specifier)
}

// resolveFile applies the resolution logic above to a filesystem,
// given a path. It's also used for the Relative resolver.
func resolveFile(base http.FileSystem, p string) ([]byte, vfs.Location, []Candidate) {
	var candidates []Candidate
	if path.Ext(p) == "" {
		candidates = append(candidates, Candidate{p + ".js", extensionRule})
		_, err := vfsutil.Stat(base, p+".js")
		switch {
		case os.IsNotExist(err):
			p = path.Join(p, "index.js")
			candidates = append(candidates, Candidate{p, indexRule})
		case err != nil:
			return nil, vfs.Nowhere, candidates
		default:
			p = p + ".js"
		}
	} else {
		candidates = append(candidates, Candidate{p, verbatimRule})
	}

	if _, err := vfsutil.Stat(base, p); err != nil {
		return nil, vfs.Nowhere, candidates
	}
	codeBytes, err := vfsutil.ReadFile(base, p)
	if err != nil {
		log.Fatal(err)
	}
	return codeBytes, vfs.Location{Vfs: base, Path: p}, candidates
}
