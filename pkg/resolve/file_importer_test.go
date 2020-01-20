package resolve

import (
	"testing"
)

func TestFileResolve(t *testing.T) {
	fr := FileImporter{vfs: ScriptBase(".").Vfs}
	_, loc, _ := fr.Import(ScriptBase("_ignored_"), "testfiles/foo", "<test>")
	if loc.Path != "testfiles/foo.js" {
		t.Errorf("did not resolve testfiles/foo, got %q", loc.Path)
	}
}

func TestVerbatimResolve(t *testing.T) {
	fr := FileImporter{vfs: ScriptBase("testfiles").Vfs}
	_, loc, _ := fr.Import(ScriptBase("_ignored_"), "foo.js", "<test>")
	if loc.Path != "foo.js" {
		t.Errorf("did not resolve foo.js, got %q", loc.Path)
	}
}

func TestIndexResolve(t *testing.T) {
	fr := FileImporter{vfs: ScriptBase("testfiles").Vfs}
	_, loc, _ := fr.Import(ScriptBase("_ignored_"), "node_modules/modfoo", "<test>")
	if loc.Path != "node_modules/modfoo/index.js" {
		t.Errorf("did not resolve foo.js, got %q", loc.Path)
	}
}
