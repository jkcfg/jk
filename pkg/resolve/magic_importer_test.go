package resolve

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkcfg/jk/pkg/vfs"
)

func TestMagicImporter(t *testing.T) {
	spec := "@jkcfg/magic"
	mod := []byte("export default foo = 1;")
	importer := MagicImporter{
		Specifier: spec,
		Generate: func(vfs.Location) ([]byte, string) {
			return mod, "foo.js"
		},
	}

	// Get the module if you 1. ask for the right thing and 2. are a std module yourself
	bytes, _, candidates := importer.Import(ScriptBase("."), spec, "@jkcfg/std/other")
	assert.Equal(t, mod, bytes)
	assert.Nil(t, candidates)

	var loc vfs.Location
	// Doesn't respond to anything that is not the magic module
	bytes, loc, _ = importer.Import(ScriptBase("."), "@jkcfg/notmagic", "@jkcfg/std/other")
	assert.Equal(t, vfs.Nowhere, loc)
	assert.Nil(t, bytes)

	// Doesn't respond if the referrer is not in std
	bytes, loc, _ = importer.Import(ScriptBase("."), spec, "rando.js")
	assert.Equal(t, vfs.Nowhere, loc)
	assert.Nil(t, bytes)

	// Make it public; should be accessible from modules outside std
	importer = MagicImporter{
		Specifier: spec,
		Generate: func(vfs.Location) ([]byte, string) {
			return mod, "foo.js"
		},
		Public: true,
	}

	bytes, _, candidates = importer.Import(ScriptBase("."), spec, "rando.js")
	assert.Equal(t, mod, bytes)
	assert.Nil(t, candidates)
}
