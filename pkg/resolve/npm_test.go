package resolve

import (
	"os"
	"testing"
)

func TestNodeModuleImport(t *testing.T) {
	// Depends on the layout of directories and files under `./testfiles`

	node := &NodeImporter{ModuleBase: "testfiles"}

	// referrer is not important, since the context is captured in
	// basePath.

	test := func(name, base, path string) {
		t.Run(name, func(t *testing.T) {
			bytes, path, candidates := node.Import(base, path, "stdin")
			if bytes == nil {
				t.Error("did not resolve", path)
				println("candidates:")
				for _, c := range candidates {
					println("  ", c.Path, c.Rule)
				}
			}
		})
	}

	// These refer to the file testfiles/foo.js
	test("exact", "testfiles", "./foo.js")
	test("sans .js", "testfiles", "./foo")

	// Using an index file
	test("index", "testfiles", "./pkg/guess/bar")

	// Via package.json
	test("via package.json with exact path", "testfiles", "./pkg/exact")
	test("via package.json and index", "testfiles", "./pkg/guess")

	// In variously nested node_modules locations (i.e., non-relative)
	test("in a node_modules package, as an index file", "testfiles/", "modfoo")
	test("in a package, in a sub-directory, then package.json", "testfiles/", "modfoo/lib")
	test("in a package, in a sub-directory, preferring file to dir", "testfiles/", "modfoo/lib/bar")

	test("under ../../node_modules", "testfiles/pkg/foo", "modfoo")
}

func TestNodeModuleDotBase(t *testing.T) {
	// This tests using a dot path '.' as the module base, so we must
	// be _in_ the directory with the test files for anything to be
	// resolvable.
	if err := os.Chdir("testfiles"); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir("..")
	node := &NodeImporter{ModuleBase: "."}

	test := func(name, base, path string) {
		t.Run(name, func(t *testing.T) {
			bytes, path, candidates := node.Import(base, path, "stdin")
			if bytes == nil {
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
	node := &NodeImporter{ModuleBase: "testfiles/pkg"}
	bytes, path, _ := node.Import("testfiles/pkg", "modfoo", "stdin")
	if bytes != nil || path != "" {
		t.Errorf("Expected failure to resolve modfoo, but found at %s", path)
	}
}
