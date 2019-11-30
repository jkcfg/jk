package resolve

import (
	"path"
	"strings"

	"github.com/jkcfg/jk/pkg/vfs"
)

func isRelative(specifier string) bool {
	return strings.HasPrefix(specifier, "./") || strings.HasPrefix(specifier, "../")
}

// Relative resolves specifiers that start with './' or '../', by
// looking for a file relative to the importing module (as given by
// the base argument).
type Relative struct {
}

// Import implements Importer with the relative import rules.
func (r Relative) Import(base vfs.Location, specifier, referrer string) ([]byte, vfs.Location, []Candidate) {
	if isRelative(specifier) {
		return resolveFile(base.Vfs, path.Join(base.Path, specifier))
	}
	return nil, vfs.Nowhere, nil
}
