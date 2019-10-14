package std

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jkcfg/jk/pkg/record"
)

// ResourceBaser is an interface for getting base paths for resources.
type ResourceBaser interface {
	ResourceBase(string) (string, bool)
}

// ReadBase resolves relative paths, and resources (module-relative
// paths). Reads outside the base are forbidden and will return an
// error.
type ReadBase struct {
	Path      string
	Resources ResourceBaser
	Recorder  *record.Recorder
}

// getPath resolves a path and an optional module reference; to an
// base path (either the input directory or the module directory), and
// a path relative to that.
func (r ReadBase) getPath(path, module string) (string, string, error) {
	base := r.Path
	if module != "" {
		modBase, ok := r.Resources.ResourceBase(module)
		if !ok {
			return "", "", fmt.Errorf("read from unknown module")
		}
		base = modBase
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(base, path)
	}

	rel, err := filepath.Rel(base, path)
	if err != nil {
		return "", "", err
	}
	if strings.HasPrefix(rel, "..") {
		return "", "", fmt.Errorf("reads outside base path forbidden")
	}

	return base, rel, nil
}
