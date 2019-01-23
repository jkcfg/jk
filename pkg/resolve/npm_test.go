package resolve

import (
	"testing"
)

func TestNodeModuleImport(t *testing.T) {
	// Depends on the layout of directories and files under `./testfiles`

	node := &NodeModulesImporter{ModuleBase: "testfiles"}

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
}
