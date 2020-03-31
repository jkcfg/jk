package std

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkcfg/jk/pkg/vfs"
)

func TestSandbox(t *testing.T) {
	mods := NewModuleResources()

	moduleBase := "/modules"
	inputDir := "testfiles"
	outputDir := "/writeroot"

	sb := Sandbox{
		Base: vfs.Location{
			Vfs:  vfs.Empty,
			Path: inputDir,
		},
		WriteRoot: outputDir,
		Modules:   mods,
	}

	// reads are relative to the module location or input dir
	readSandboxSucceeds := func(t *testing.T, module, base string) {
		t.Run("read from subdir", func(t *testing.T) {
			loc, err := sb.getReadPath("./subdir/foo.json", module)
			assert.NoError(t, err)
			assert.Equal(t, base+"/subdir/foo.json", loc.Path)
		})
		t.Run("read with internal parent path", func(t *testing.T) {
			loc, err := sb.getReadPath("foo/bar/../index.yaml", module)
			assert.NoError(t, err)
			assert.Equal(t, base+"/foo/index.yaml", loc.Path)
		})
	}

	readOutsideSucceeds := func(t *testing.T, module string) {
		t.Run("read resource absolute path", func(t *testing.T) {
			loc, err := sb.getReadPath("/root/foo.json", module)
			assert.NoError(t, err)
			assert.Equal(t, "/root/foo.json", loc.Path)
		})
		t.Run("read resource parent path", func(t *testing.T) {
			loc, err := sb.getReadPath("../foo.json", module)
			assert.NoError(t, err)
			// NB both moduleBase and inputDir are one element, so a
			// parent path ends up at the root.
			assert.Equal(t, "/foo.json", loc.Path)
		})
	}

	readOutsideFails := func(t *testing.T, module string) {
		t.Run("fail read with absolute path", func(t *testing.T) {
			loc, err := sb.getReadPath("/tmp/index.yaml", module)
			assert.Error(t, err)
			assert.Equal(t, vfs.Nowhere, loc)
		})
		t.Run("fail read with parent path", func(t *testing.T) {
			loc, err := sb.getReadPath("sub/../../index.yaml", module)
			assert.Error(t, err)
			assert.Equal(t, vfs.Nowhere, loc)
		})
	}

	// writes are always to the host filesystem, and relative to the
	// output directory.
	writeSandboxSucceeds := func(t *testing.T, module string) {
		t.Run("allow write to subdir", func(t *testing.T) {
			path, err := sb.getWritePath("sub/foo.json", module)
			assert.NoError(t, err)
			assert.Equal(t, outputDir+"/sub/foo.json", path)
		})
	}

	writeOutsideSucceeds := func(t *testing.T, module string) {
		t.Run("write to absolute path", func(t *testing.T) {
			path, err := sb.getWritePath("/tmp/alternate/foo.json", module)
			assert.NoError(t, err)
			assert.Equal(t, "/tmp/alternate/foo.json", path)
		})
		t.Run("write to parent path", func(t *testing.T) {
			path, err := sb.getWritePath("../foo.json", module)
			assert.NoError(t, err)
			// NB relies on the output dir being just one path element
			assert.Equal(t, "/foo.json", path)
		})
	}

	writeSandboxFails := func(t *testing.T, module string) {
		t.Run("fail write to module path", func(t *testing.T) {
			path, err := sb.getWritePath("sub/foo.json", module)
			assert.Error(t, err)
			assert.Equal(t, "", path)
		})
	}

	writeOutsideFails := func(t *testing.T, module string) {
		t.Run("fail writes to parent path", func(t *testing.T) {
			path, err := sb.getWritePath("./foo/../../bar.yaml", module)
			assert.Error(t, err)
			assert.Equal(t, "", path)
		})
		t.Run("fail writes to absolute path", func(t *testing.T) {
			path, err := sb.getWritePath("/tmp/bar.yaml", module)
			assert.Error(t, err)
			assert.Equal(t, "", path)
		})
	}

	// ----------

	t.Run("without module", func(t *testing.T) {
		readSandboxSucceeds(t, "", inputDir)
		writeSandboxSucceeds(t, "")

		readOutsideFails(t, "")
		writeOutsideFails(t, "")
	})

	t.Run("readonly, sandboxed module", func(t *testing.T) {
		// Read-only access to a module
		readonlySandboxMod := ModuleAccess{
			Loc: vfs.Location{
				Vfs:  vfs.Empty,
				Path: moduleBase,
			},
		}
		readonlyToken := mods.registerModuleAccess(readonlySandboxMod)

		readSandboxSucceeds(t, readonlyToken, moduleBase)

		writeSandboxFails(t, readonlyToken)
		readOutsideFails(t, readonlyToken)
		writeOutsideFails(t, readonlyToken)
	})

	t.Run("readonly, non-sandboxed module", func(t *testing.T) {
		// Still readonly, but can read parent/absolute paths
		readonlyEscapeMod := ModuleAccess{
			Loc: vfs.Location{
				Vfs:  vfs.Empty,
				Path: moduleBase,
			},
			AllowPathsOutsideSandbox: true,
		}
		readonlyEscapeToken := mods.registerModuleAccess(readonlyEscapeMod)

		readSandboxSucceeds(t, readonlyEscapeToken, moduleBase)
		readOutsideSucceeds(t, readonlyEscapeToken)

		writeSandboxFails(t, readonlyEscapeToken)
		writeOutsideFails(t, readonlyEscapeToken)
	})

	t.Run("writable, sandboxed module", func(t *testing.T) {
		// Can write to module and outputDir, but not outside sandbox
		writableMod := ModuleAccess{
			Loc: vfs.Location{
				Vfs:  vfs.Empty,
				Path: moduleBase,
			},
			AllowWriteToHost: true,
		}
		writableToken := mods.registerModuleAccess(writableMod)

		readSandboxSucceeds(t, writableToken, moduleBase)
		writeSandboxSucceeds(t, writableToken)

		readOutsideFails(t, writableToken)
		writeOutsideFails(t, writableToken)
	})

	t.Run("writable, non-sandboxed module", func(t *testing.T) {
		// Can read and write anywhere.
		writableEscapeMod := ModuleAccess{
			Loc: vfs.Location{
				Vfs:  vfs.Empty,
				Path: moduleBase,
			},
			AllowPathsOutsideSandbox: true,
			AllowWriteToHost:         true,
		}
		writableEscapeToken := mods.registerModuleAccess(writableEscapeMod)

		readSandboxSucceeds(t, writableEscapeToken, moduleBase)
		writeSandboxSucceeds(t, writableEscapeToken)
		readOutsideSucceeds(t, writableEscapeToken)
		writeOutsideSucceeds(t, writableEscapeToken)
	})

}
