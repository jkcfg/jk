package resolve

import (
	"testing"
)

func TestRelativeResolve(t *testing.T) {
	relative := Relative{}
	base := ScriptBase("testfiles")

	test := func(name, path string) {
		t.Run(name, func(t *testing.T) {
			bytes, loc, candidates := relative.Import(base, path, "stdin")
			if bytes == nil || loc.Path == "" {
				t.Error("did not resolve", path)
				println("candidates:")
				for _, c := range candidates {
					println("  ", c.Path, c.Rule)
				}
			}
		})
	}

	// These refer to the file testfiles/foo.js
	test("exact", "./foo.js")
	test("sans .js", "./foo")
	// Using an index file
	test("index", "./pkg/guess/bar")
}
