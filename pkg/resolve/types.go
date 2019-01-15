package resolve

import (
	v8 "github.com/ry/v8worker2"
)

// Loader is an object able to load a ES 2015 module.
type Loader interface {
	LoadModule(scriptName string, code string, resolve v8.ModuleResolverCallback) error
}

// Importer is a object resolving a import to actual JS code.
type Importer interface {
	// Resolve a specifier (e.g., `my-module/foo') to a specific path
	// and file contents. Also returns a list of the interpretations
	// of the specifier attempted, including that returned.
	Import(basePath, specifier, referrer string) (data []byte, path string, candidates []string)
}
