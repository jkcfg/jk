package resolve

import (
	"net/http"
	"testing"
)

func TestModuleResolve(t *testing.T) {
	fr := FileImporter{vfs: http.Dir(".")}
	_, loc, _ := fr.Import(ScriptBase("."), "testfiles/foo.js", "<test>")
	if loc.Path != "testfiles/foo.js" {
		t.Errorf("did not resolve testfiles/foo, got %q", loc.Path)
	}
}
