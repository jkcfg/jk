package std

import (
	"fmt"
	"path"
	"strings"

	"github.com/jkcfg/jk/pkg/record"
	"github.com/jkcfg/jk/pkg/vfs"
)

// ResourceBaser is an interface for getting base paths for resources.
type ResourceBaser interface {
	ResourceBase(string) (vfs.Location, bool)
}

// ReadBase resolves relative paths, and resources (module-relative
// paths). Reads outside the base are forbidden and will return an
// error.
type ReadBase struct {
	Base      vfs.Location
	Resources ResourceBaser
	Recorder  *record.Recorder
}

// getPath resolves a path and an optional module reference; to an
// base location (either the input directory or the module directory),
// and a path relative to that. For some uses we want to know the
// relative path, so it's kept separate rather than returning just the
// fully-resolved location.
func (r ReadBase) getPath(p, module string) (vfs.Location, string, error) {
	// Absolute paths are never allowed
	if path.IsAbs(p) {
		return vfs.Nowhere, "", fmt.Errorf("absolute paths are forbidden")
	}

	// Paths outside the input directory (Base) or module directory
	// are considered forbidden.
	//
	// path.Clean will bring any parent paths (`..`) to the beginning
	// of the path. Anything that begins with a parent path is
	// forbidden.
	//
	// Note that it's possible to have a cleaned path that is _within_
	// the input directory, but begins with a parent path (e.g., if
	// you name the directory you're in). But it's not possible to
	// have a relative path that escapes the input directory _without_
	// the cleaned version starting with a parent path. I.e.,
	//
	//     invalid path -> starts with `..`
	//
	// but not vice versa. So, this will rule out some technically OK
	// paths, but these are always able to be expressed such that they
	// will be allowed (e.g., instead of `../foo/bar.yaml` while
	// assuming that CWD is `foo`, use `./bar.yaml`).
	relPath := path.Clean(p)
	if strings.HasPrefix(relPath, "..") {
		return vfs.Nowhere, "", fmt.Errorf("reading from a parent path is forbidden")
	}

	base := r.Base
	if module != "" {
		modBase, ok := r.Resources.ResourceBase(module)
		if !ok {
			return vfs.Nowhere, "", fmt.Errorf("read from unknown module")
		}
		base = modBase
	}

	return base, relPath, nil
}
