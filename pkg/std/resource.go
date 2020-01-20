package std

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/jkcfg/jk/pkg/vfs"
)

// ModuleResources keeps track of the base paths for modules, as well
// as generating the magic modules when they are imported.
type ModuleResources struct {
	mu sync.RWMutex
	// module hash -> Location for resource reads
	modules map[string]vfs.Location
	salt    []byte
}

// NewModuleResources initialises a new ModuleResources
func NewModuleResources() *ModuleResources {
	r := &ModuleResources{
		modules: map[string]vfs.Location{},
	}
	r.salt = make([]byte, 32)
	rand.Read(r.salt)
	return r
}

// ResourceBase provides the module base location given the hash.
func (r *ModuleResources) ResourceBase(hash string) (vfs.Location, bool) {
	r.mu.RLock()
	loc, ok := r.modules[hash]
	r.mu.RUnlock()
	return loc, ok
}

// MakeModule generates resource module code (and path) given the
// importing module's base path.
func (r *ModuleResources) MakeModule(base vfs.Location) ([]byte, string) {
	hash := sha256.New()
	hash.Write([]byte(base.Path)) // TODO needs more'n this, not unique enough
	hash.Write(r.salt)
	moduleHash := fmt.Sprintf("%x", hash.Sum(nil))
	r.mu.Lock()
	r.modules[moduleHash] = base
	r.mu.Unlock()

	code := `
import * as std from '@jkcfg/std';
import * as fs from '@jkcfg/std/fs';

const module = %q;

function read(path, {...rest} = {}) {
  return std.read(path, {...rest, module});
}

function dir(path) {
  return fs.dir(path, { module });
}

function info(path) {
  return fs.info(path, { module });
}

function withModuleRef(fn) {
  return fn(module);
};

export { read, dir, info, withModuleRef };
`
	return []byte(fmt.Sprintf(code, moduleHash)), "resource:" + base.Path // TODO this one too, meant to be unique
}
