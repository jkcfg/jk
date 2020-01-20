package resolve

import (
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
//   - `/foo/bar.{js,mjs}`
//   - `/foo/bar/index.{js,mjs}
type FileImporter struct {
	vfs vfs.FileSystem
}

// NewFileImporter constructs a FileImport given a filesystem
func NewFileImporter(vfs vfs.FileSystem) *FileImporter {
	return &FileImporter{vfs: vfs}
}

const (
	extensionRulePrefix = "<path> -> <path>"
	indexRulePrefix     = "<path> -> <path>/index"
	verbatimRule        = "verbatim"
)

var moduleExtensions = []string{".mjs", ".js"}

// Import implements importer. Note that the file import only ever
// cares to look in the root of the given filesystem. It doesn't care
// where the importing module is located.
func (fi *FileImporter) Import(base vfs.Location, specifier, referrer string) ([]byte, vfs.Location, []Candidate) {
	if isRelative(specifier) {
		return nil, vfs.Nowhere, nil
	}
	return resolvePath(fi.vfs, specifier)
}

// resolvePath tries the rules as given above
func resolvePath(fs vfs.FileSystem, p string) ([]byte, vfs.Location, []Candidate) {
	bytes, loc, fileCandidates := resolveFile(fs, p)
	if bytes != nil {
		return bytes, loc, fileCandidates
	}
	bytes, loc, indexCandidates := resolveIndex(fs, p)
	return bytes, loc, append(fileCandidates, indexCandidates...)
}

// resolveFile tries to load a path as though it referred to a
// file. No bytes returned means failure.
func resolveFile(fs vfs.FileSystem, p string) ([]byte, vfs.Location, []Candidate) {
	candidates := []Candidate{{p, verbatimRule}}
	bytes, err := vfsutil.ReadFile(fs, p)
	if err == nil {
		return bytes, vfs.Location{Vfs: fs, Path: p}, candidates
	}

	bytes, loc, extCandidates := resolveGuessedFile(fs, p)
	return bytes, loc, append(candidates, extCandidates...)
}

// resolveGuessedFile tries to apply the rule `specifier ->
// specifier.{js,mjs} to resolve a path p within filesystem fs.
func resolveGuessedFile(fs vfs.FileSystem, p string) ([]byte, vfs.Location, []Candidate) {
	var candidates []Candidate
	for _, ext := range moduleExtensions {
		f := p + ext
		candidates = append(candidates, Candidate{f, extensionRulePrefix + ext})
		bytes, err := vfsutil.ReadFile(fs, f)
		if err == nil {
			return bytes, vfs.Location{Vfs: fs, Path: f}, candidates
		}
	}
	return nil, vfs.Nowhere, candidates
}

// resolveIndex tries to load the default index files, assuming the path
// is a directory.
func resolveIndex(fs vfs.FileSystem, base string) ([]byte, vfs.Location, []Candidate) {
	var candidates []Candidate
	for _, ext := range moduleExtensions {
		p := path.Join(base, "index"+ext)
		candidates = append(candidates, Candidate{p, indexRulePrefix + ext})
		bytes, err := vfsutil.ReadFile(fs, p)
		if err == nil {
			return bytes, vfs.Location{Vfs: fs, Path: p}, candidates
		}
	}
	return nil, vfs.Nowhere, candidates
}

// // resolveFile applies the resolution logic above to a filesystem,
// // given a path. It's also used for the Relative resolver.
// func resolveFile(base http.FileSystem, p string) ([]byte, vfs.Location, []Candidate) {
// 	var candidates []Candidate
// 	if path.Ext(p) == "" {
// 		candidates = append(candidates, Candidate{p + ".js", extensionRule})
// 		_, err := vfsutil.Stat(base, p+".js")
// 		switch {
// 		case os.IsNotExist(err):
// 			p = path.Join(p, "index.js")
// 			candidates = append(candidates, Candidate{p, indexRule})
// 		case err != nil:
// 			return nil, vfs.Nowhere, candidates
// 		default:
// 			p = p + ".js"
// 		}
// 	} else {
// 		candidates = append(candidates, Candidate{p, verbatimRule})
// 	}

// 	if _, err := vfsutil.Stat(base, p); err != nil {
// 		return nil, vfs.Nowhere, candidates
// 	}
// 	codeBytes, err := vfsutil.ReadFile(base, p)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return codeBytes, vfs.Location{Vfs: base, Path: p}, candidates
// }
