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
	// The location in a virtual filesystem that is the top-most that
	// can be read from
	Base vfs.Location
	// The top-most directory that can be written to
	WriteRoot string
	// Look-up for resources (i.e., for module-relative reads)
	Resources ResourceBaser
	// For recording each read or write
	Recorder *record.Recorder
}

// getReadPath resolves a path and an optional module reference, to a
// location for reading.
func (s Sandbox) getReadPath(p, module string) (vfs.Location, error) {
	base := s.Base
	if module != "" {
		modBase, ok := s.Resources.ResourceBase(module)
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

// getWritePath verifies the path given and resolves it relative to
// the output directory.
func (s Sandbox) getWritePath(p string) (string, error) {
	if p == "" {
		return p, nil
	}
	if path.IsAbs(p) {
		return "", fmt.Errorf("writing to an absolute path is forbidden")
	}

	p = path.Clean(p)
	// See note in `getReadPath` about parent paths
	if strings.HasPrefix(p, "../") {
		return "", fmt.Errorf("writing to parent path is forbidden")
	}
	p = path.Join(s.WriteRoot, p)
	return p, nil
}
