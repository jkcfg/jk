package resolve

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/jkcfg/jk/pkg/vfs"
	"github.com/shurcooL/httpfs/vfsutil"
)

/* ## Module resolution for npm packages

[npm](https://docs.npmjs.com/about-npm/) is a convenient way to
distribute JavaScript modules. By and large these are for Node.JS (and
use CommonJS modules) or for browsers (and have whatever unholy
packaging was considered best practice at the time); but, it's
possible to publish ES2015 modules as well, either consisting of just
ES6 code, or hybrid packages that can be used with more than one
module system.

Here's how npm packages work:

 1. When you `npm install foo`, npm creates a directory
 `./node_modules/foo` with the contents of the package in it.

 2. In that directory, there will be a file `package.json`, which has
 various metadata, and some fields describing how `foo` is to be used;

 3. If the `package.json` has entries under `.dependencies`, these
 will also be downloaded by npm. It may put them inside the foo
 directory (i.e., under `./node_modules/foo/node_modules`), or flatten
 them out into the directory that `foo` is in (i.e., under
 `./node_modules`). Newer npm releases tend to try and flatten
 dependencies.

 4. If the `package.json` file has an field `.main`, that points to
 the file to be used when the package is referred to by its directory,
 e.g., if you import `foo`.

The algorithm for resolving a module path is described in
<https://nodejs.org/api/modules.html#modules_all_together>. This file
adapts that algorithm to account for only supporting ES2015 modules,
by

 - using the file extensions `.js`, `.mjs`
 - using the field `module` in `package.json`, instead of `main`
 - not looking for node_modules directories above the directory given
   as the top-level.

*/

// NodeImporter is an implementation of Importer that uses a resolution
// algorithm adapted from Node.JS's, in order to support modules installed
// with npm.
type NodeImporter struct {
	vfs http.FileSystem
}

// NewNodeImporter constructs a NodeImporter using the given filesystem
func NewNodeImporter(vfs http.FileSystem) *NodeImporter {
	return &NodeImporter{vfs: vfs}
}

func (n *NodeImporter) locationAt(path string) vfs.Location {
	return vfs.Location{Vfs: n.vfs, Path: path}
}

// Import is the entry point into the module resolution algorithm.
func (n *NodeImporter) Import(base vfs.Location, specifier, referrer string) ([]byte, vfs.Location, []Candidate) {
	if path.IsAbs(specifier) {
		log.Fatalf("absolute import path %q not allowed in %q", specifier, referrer)
	}
	if isRelative(specifier) {
		return nil, vfs.Nowhere, nil
	}
	return n.loadAsModule(specifier, base.Path)
}

var moduleExtensions = []string{".mjs", ".js"}

// loadAsFile tries to load a path as though it referred to a file. No
// bytes returned means failure.
func (n *NodeImporter) loadAsFile(path string) ([]byte, vfs.Location, []Candidate) {
	candidates := []Candidate{{path, verbatimRule}}
	bytes, err := vfsutil.ReadFile(n.vfs, path)
	if err == nil {
		return bytes, n.locationAt(path), candidates
	}

	bytes, loc, extCandidates := n.loadGuessedFile(path)
	return bytes, loc, append(candidates, extCandidates...)
}

func (n *NodeImporter) loadGuessedFile(path string) ([]byte, vfs.Location, []Candidate) {
	var candidates []Candidate
	for _, ext := range moduleExtensions {
		p := path + ext
		candidates = append(candidates, Candidate{p, extensionRulePrefix + ext})
		bytes, err := vfsutil.ReadFile(n.vfs, p)
		if err == nil {
			return bytes, n.locationAt(p), candidates
		}
	}
	return nil, vfs.Nowhere, candidates
}

// loadAsPath attempts to load a path when it's unknown whether it
// refers to a file or a directory.
func (n *NodeImporter) loadAsPath(path string) ([]byte, vfs.Location, []Candidate) {
	bytes, loc, fileCandidates := n.loadAsFile(path)
	if bytes != nil {
		return bytes, loc, fileCandidates
	}

	info, err := vfsutil.Stat(n.vfs, path)
	switch {
	case os.IsNotExist(err):
		// loadAsFile will already have included the possibility of
		// the path as-is as a candidate
		return nil, vfs.Nowhere, fileCandidates
	case err != nil:
		log.Fatal(err)

	case info.IsDir():
		bytes, loc, dirCandidates := n.loadAsDir(path)
		return bytes, loc, append(fileCandidates, dirCandidates...)
	}
	return nil, vfs.Nowhere, fileCandidates
}

// loadIndex tries to load the default index files, assuming the path
// is a directory.
func (n *NodeImporter) loadIndex(base string) ([]byte, vfs.Location, []Candidate) {
	var candidates []Candidate
	for _, ext := range moduleExtensions {
		p := path.Join(base, "index"+ext)
		candidates = append(candidates, Candidate{p, indexRulePrefix + ext})
		bytes, err := vfsutil.ReadFile(n.vfs, p)
		if err == nil {
			return bytes, n.locationAt(p), candidates
		}
	}
	return nil, vfs.Nowhere, candidates
}

// loadAsDir attempts to load a path which is known to be a directory.
func (n *NodeImporter) loadAsDir(dir string) ([]byte, vfs.Location, []Candidate) {
	var candidates []Candidate

	packageJSONPath := path.Join(dir, "package.json")
	packageJSON, _ := vfsutil.ReadFile(n.vfs, packageJSONPath)
	if packageJSON != nil {
		var pkg struct{ Module string }
		if err := json.Unmarshal(packageJSON, &pkg); err == nil && pkg.Module != "" {
			module := path.Join(dir, pkg.Module)
			// .module is treated as through it were a file (but not a directory)
			bytes, loc, pkgCandidates := n.loadAsFile(module)
			// TODO(michael) consider transformating these candidates
			// (and below) to reflect the indirection through
			// package.json
			qualifyCandidates(pkgCandidates, "via .module in "+packageJSONPath)
			candidates = append(candidates, pkgCandidates...)
			if bytes != nil {
				return bytes, loc, candidates
			}
			// .. or a directory with an index (but not another package.json)
			bytes, loc, modIndexCandidates := n.loadIndex(module)
			qualifyCandidates(modIndexCandidates, "via .module in "+packageJSONPath)
			candidates = append(candidates, modIndexCandidates...)
			if bytes != nil {
				return bytes, loc, candidates
			}
		}
	}
	bytes, loc, indexCandidates := n.loadIndex(dir)
	return bytes, loc, append(candidates, indexCandidates...)
}

// loadAsModule attempts to load a specifier as though it referred to
// a package in (potentially nested) node_modules directories.
func (n *NodeImporter) loadAsModule(specifier, base string) (_ []byte, _ vfs.Location, candidates []Candidate) {
	defer func() {
		qualifyCandidates(candidates, "via NPM resolution")
	}()

	bits := strings.Split(base, "/")
	for i := len(bits); i >= 0; i-- {
		if i > 0 && bits[i-1] == "node_modules" {
			continue
		}
		path := path.Join(append(bits[:i], "node_modules", specifier)...)
		bytes, loc, pathCandidates := n.loadAsPath(path)
		candidates = append(candidates, pathCandidates...)
		if bytes != nil {
			return bytes, loc, candidates
		}
	}
	return nil, vfs.Nowhere, candidates
}

func qualifyCandidates(cs []Candidate, extra string) {
	for i := range cs {
		cs[i].Rule = cs[i].Rule + ", " + extra
	}
}
