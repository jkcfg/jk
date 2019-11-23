package resolve

import (
	"testing"
)

func TestResolve(t *testing.T) {
	fr := FileImporter{ModuleBase: "testfiles"}
	_, path, _ := fr.Import(".", "./testfiles/foo", "<test>")
	if path != "testfiles/foo.js" {
		t.Errorf("did not resolve ./testfiles/foo, got %q", path)
	}
}

func TestOutsideModuleBase(t *testing.T) {
	fr := FileImporter{ModuleBase: "testfiles/pkg"}
	_, path, _ := fr.Import(".", "testfiles/foo.js", "<test>")
	if path != "" {
		t.Errorf("should not resolve; but was resolved to %q", path)
	}
}
