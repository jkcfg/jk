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
		"std",
		"@jkcfg/std", "test.js", "/path/to/dir",
		true, "@jkcfg/std/std.js",
	}, {
		"std.js",
		"@jkcfg/std/std.js", "test.js", "/path/to/dir",
		true, "@jkcfg/std/std.js",
	}, {
		"std/param",
		"@jkcfg/std/param", "test.js", "/path/to/dir",
		true, "@jkcfg/std/std_param.js",
	}, {
		"std/param.js",
		"@jkcfg/std/param.js", "test.js", "/path/to/dir",
		true, "@jkcfg/std/std_param.js",
	}, {
		// Users cannot load non-exported modules.
		"not-public",
		"@jkcfg/std/foo", "test.js", "/path/to/dir",
		false, "",
	}, {
		// We can still import std modules from the std code itself.
		"internal",
		"std_log", "@jkcfg/std/std.js", "@jkcfg/std",
		true, "@jkcfg/std/std_log.js",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := &StdImporter{
				PublicModules: []StdPublicModule{{
					ExternalName: "std.js", InternalModule: "std.js",
				}, {
					ExternalName: "param.js", InternalModule: "std_param.js",
				}},
			}

			source, resolved, _ := i.Import(test.base, test.specifier, test.referrer)
			if !test.valid {
				assert.Equal(t, 0, len(source))
				return
			}
			assert.NotEqual(t, 0, len(source))
			assert.Equal(t, test.resolved, resolved)
		})
	}
}
