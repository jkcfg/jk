package resolve

import (
	v8 "github.com/ry/v8worker2"
)

// Loader is an object able to load a ES 2015 module.
type Loader interface {
	LoadModule(scriptName string, code string, resolve v8.ModuleResolverCallback) error
}
