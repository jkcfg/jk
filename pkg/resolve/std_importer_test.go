package resolve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdImporter(t *testing.T) {
	tests := []struct {
		name                string
		specifier, referrer string
		base                string

		valid    bool
		resolved string
	}{{
		"default",
		"@jkcfg/std", "test.js", "/path/to/dir",
		true, "@jkcfg/std/index.js",
	}, {
		"index.js",
		"@jkcfg/std/index.js", "test.js", "/path/to/dir",
		true, "@jkcfg/std/index.js",
	}, {
		"std/param",
		"@jkcfg/std/param", "test.js", "/path/to/dir",
		true, "@jkcfg/std/param.js",
	}, {
		"std/param.js",
		"@jkcfg/std/param.js", "test.js", "/path/to/dir",
		true, "@jkcfg/std/param.js",
	}, {
		// Users cannot load non-exported modules.
		"not-public",
		"@jkcfg/std/foo", "test.js", "/path/to/dir",
		false, "",
	}, {
		// We can still import std modules from the std code itself.
		"internal",
		"./log", "@jkcfg/std/index.js", "@jkcfg/std",
		true, "@jkcfg/std/log.js",
	}, {
		"internal-relative-deep-path",
		"./flatbuffers", "@jkcfg/std/internal/deferred.js", "@jkcfg/std/internal",
		true, "@jkcfg/std/internal/flatbuffers.js",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := &StdImporter{
				PublicModules: []string{"index.js", "param.js"},
			}

			source, resolved, _ := i.Import(test.base, test.specifier, test.referrer)
			if !test.valid {
				assert.Empty(t, source)
				return
			}
			assert.NotEmpty(t, source)
			assert.Equal(t, test.resolved, resolved)
		})
	}
}
