package std

import (
	"fmt"
	"path"
	"path/filepath"
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
// and a path relative to that.
func (r ReadBase) getPath(p, module string) (vfs.Location, string, error) {
	base := r.Base
	if module != "" {
		modBase, ok := r.Resources.ResourceBase(module)
		if !ok {
			return vfs.Nowhere, "", fmt.Errorf("read from unknown module")
		}
		base = modBase
	}

	if !path.IsAbs(p) {
		p = path.Join(base.Path, p)
	}

	rel, err := filepath.Rel(base.Path, p) // TODO no Rel in path
	if err != nil {
		return vfs.Nowhere, "", err
	}
	if strings.HasPrefix(rel, "..") {
		return vfs.Nowhere, "", fmt.Errorf("reads outside base path forbidden")
	}

	return base, rel, nil
}
