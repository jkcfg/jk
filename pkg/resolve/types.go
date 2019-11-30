package resolve

import (
	v8 "github.com/jkcfg/v8worker2"

	"github.com/jkcfg/jk/pkg/vfs"
)

// Loader is an object able to load a ES 2015 module.
type Loader interface {
	LoadModule(scriptName string, code string, resolve v8.ModuleResolverCallback) error
}

// Candidate is a path that was considered when resolving an import,
// and the explanation (resolution rule) for why it was considered.
type Candidate struct {
	Path string
	Rule string
}

// Importer is a object resolving a import to actual JS code.
type Importer interface {
	// Resolve a specifier (e.g., `my-module/foo') to a specific path
	// and file contents. Also returns a list of the interpretations
	// of the specifier attempted, including that returned.
	Import(base vfs.Location, specifier, referrer string) (data []byte, resolved vfs.Location, candidates []Candidate)
}
