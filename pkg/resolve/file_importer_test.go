package resolve

import (
	"testing"
)

func TestModuleResolve(t *testing.T) {
	fr := FileImporter{vfs: ScriptBase(".").Vfs}
	_, loc, _ := fr.Import(ScriptBase("."), "testfiles/foo.js", "<test>")
	if loc.Path != "testfiles/foo.js" {
		t.Errorf("did not resolve testfiles/foo, got %q", loc.Path)
	}
}
