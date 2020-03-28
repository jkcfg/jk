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

// Sandbox mediates access to the filesystem by resolving relative
// paths to the host filesystem and resources (module-relative
// paths). Reads outside the base are forbidden and will return an
// error.
type Sandbox struct {
	Base      vfs.Location
	Resources ResourceBaser
	Recorder  *record.Recorder
}

// getPath resolves a path and an optional module reference, to a
// location.
func (r Sandbox) getPath(p, module string) (vfs.Location, error) {
	base := r.Base
	if module != "" {
		modBase, ok := r.Resources.ResourceBase(module)
		if !ok {
			return vfs.Nowhere, fmt.Errorf("read from unknown module")
		}
		base = modBase
	}

	p = path.Clean(p)

	// If this particular base location allows parent paths, we're
	// done.
	if base.AllowParentPaths {
		if !path.IsAbs(p) {
			p = path.Join(base.Path, p)
		}
		return vfs.Location{
			Vfs:              base.Vfs,
			Path:             p,
			AllowParentPaths: true,
		}, nil
	}

	// But usually, paths outside the input directory (Base) or module
	// directory are considered forbidden.

	// Absolute paths are easy to detect
	if path.IsAbs(p) {
		return vfs.Nowhere, fmt.Errorf("reading absolute paths is forbidden")
	}

	// path.Clean as done above will bring any parent paths (`..`) to
	// the beginning of the path. Anything that begins with a parent
	// path is forbidden.
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
	if strings.HasPrefix(p, "..") {
		return vfs.Nowhere, fmt.Errorf("reading from a parent path is forbidden")
	}

	return vfs.Location{
		Vfs:  base.Vfs,
		Path: path.Join(base.Path, p),
	}, nil
}
