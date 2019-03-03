package std

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// ModuleResources keeps track of the base paths for modules, as well
// as generating the magic modules when they are imported.
type ModuleResources struct {
	// module hash -> basePath for resource reads
	modules map[string]string
	salt    []byte
}

// NewModuleResources initialises a new ModuleResources
func NewModuleResources() *ModuleResources {
	r := &ModuleResources{
		modules: map[string]string{},
	}
	r.salt = make([]byte, 32)
	rand.Read(r.salt)
	return r
}

// ResourceBase provides the module base path given the hash.
func (r *ModuleResources) ResourceBase(hash string) (string, bool) {
	path, ok := r.modules[hash]
	return path, ok
}

// MakeModule generates resource module code (and path) given the
// importing module's base path.
func (r *ModuleResources) MakeModule(basePath string) ([]byte, string) {
	hash := sha256.New()
	hash.Write([]byte(basePath))
	hash.Write(r.salt)
	moduleHash := fmt.Sprintf("%x", hash.Sum(nil))
	r.modules[moduleHash] = basePath

	code := `
import std from '@jkcfg/std';

const module = %q;

function resource(path, {...rest} = {}) {
  return std.read(path, {...rest, module});
}

export default resource;
`
	return []byte(fmt.Sprintf(code, moduleHash)), "resource:" + basePath
}
