package resolve

import (
	"os"
	"testing"
)

func TestNodeModuleImport(t *testing.T) {
	// Depends on the layout of directories and files under `./testfiles`

	node := NewNodeImporter(ScriptBase("testfiles").Vfs)

	// referrer is not important, since the context is captured in
	// basePath.

	test := func(name, base, path string) {
		t.Run(name, func(t *testing.T) {
			bytes, path, candidates := node.Import(ScriptBase(base), path, "stdin")
			if bytes == nil {
				t.Error("did not resolve", path)
				println("candidates:")
				for _, c := range candidates {
					println("  ", c.Path, c.Rule)
				}
			}
		})
	}

	// In variously nested node_modules locations (i.e., non-relative)
	test("in a node_modules package, as an index file", "/", "modfoo")
	test("in a package, in a sub-directory, then package.json", "/", "modfoo/lib")
	test("in a package, in a sub-directory, preferring file to dir", "/", "modfoo/lib/bar")

	test("under ../../node_modules", "/pkg/foo", "modfoo")
}

func TestNodeModuleDotBase(t *testing.T) {
	// This tests using a dot path '.' as the module base, so we must
	// be _in_ the directory with the test files for anything to be
	// resolvable.
	if err := os.Chdir("testfiles"); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir("..")
	node := NewNodeImporter(ScriptBase(".").Vfs)

	test := func(name, base, path string) {
		t.Run(name, func(t *testing.T) {
			bytes, loc, candidates := node.Import(ScriptBase(base), path, "stdin")
			if bytes == nil || loc.Path == "" {
				t.Error("did not resolve", path)
				println("candidates:")
				for _, c := range candidates {
					println("  ", c.Path, c.Rule)
				}
			}
		})
	}

	test("under ./node_modules, via subdir", "subdir/", "modfoo")
}

func TestNodeModuleFail(t *testing.T) {
	// Test that we _can't resolve a module that's not under the ModuleBase
	node := NewNodeImporter(ScriptBase("testfiles/pkg").Vfs)
	bytes, loc, _ := node.Import(ScriptBase("testfiles/pkg"), "modfoo", "stdin")
	if bytes != nil || loc.Path != "" {
		t.Errorf("Expected failure to resolve modfoo, but found at %s", loc.Path)
	}
}
