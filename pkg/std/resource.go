package std

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"

	"github.com/jkcfg/jk/pkg/vfs"
)

// ModuleAccess represents access to a module's resources
type ModuleAccess struct {
	// Does this module allow paths outside the sandbox?
	AllowPathsOutsideSandbox bool
	// The location within a virtual filesystem; usually, but not
	// necessarily, the root of the virtual filesystem.
	Loc vfs.Location
	// Writes don't go through a virtual filesystem. This flag means
	// "allow writes to the host filesystem with this module".
	AllowWriteToHost bool
}

// ModuleResources keeps track of the base paths for modules, as well
// as generating the magic modules when they are imported.
type ModuleResources struct {
	mu sync.RWMutex
	// module hash -> Location for resource reads
	modules map[string]ModuleAccess
	salt    []byte
}

// NewModuleResources initialises a new ModuleResources
func NewModuleResources() *ModuleResources {
	r := &ModuleResources{
		modules: map[string]ModuleAccess{},
	}
	r.salt = make([]byte, 32)
	rand.Read(r.salt)
	return r
}

// GetModuleAccess looks up a module given a token.
func (r *ModuleResources) GetModuleAccess(token string) (ModuleAccess, bool) {
	r.mu.RLock()
	mod, ok := r.modules[token]
	r.mu.RUnlock()
	return mod, ok
}

func (r *ModuleResources) registerModuleAccess(mod ModuleAccess) string {
	canonicalPath := mod.Loc.CanonicalPath()
	hash := sha256.New()
	hash.Write([]byte(canonicalPath))
	hash.Write(r.salt)
	moduleHash := fmt.Sprintf("%x", hash.Sum(nil))
	r.mu.Lock()
	r.modules[moduleHash] = mod
	r.mu.Unlock()
	return moduleHash
}

// MakeResourceModule generates resource module code (and path) given the
// importing module's base path.
func (r *ModuleResources) MakeResourceModule(mod ModuleAccess) ([]byte, string) {
	moduleHash := r.registerModuleAccess(mod)
	exports := []string{
		"read", "dir", "info", "withModuleRef",
	}
	if mod.AllowWriteToHost {
		exports = append(exports, "write")
	}

	code := `
import * as std from '@jkcfg/std';
import * as fs from '@jkcfg/std/fs';

const module = %q;

function read(path, {...rest} = {}) {
  return std.read(path, {...rest, module});
}

function write(value, path, {...rest} = {}) {
  return std.write(value, path, {...rest, module});
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

export { %s };
`
	moduleCode := fmt.Sprintf(code, moduleHash, strings.Join(exports, ", "))
	return []byte(moduleCode), "resource:" + mod.Loc.CanonicalPath()
}
