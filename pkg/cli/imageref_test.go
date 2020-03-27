package cli

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestImageRefFlag(t *testing.T) {
	positives := []string{
		"foo:v1",
		"alpine:3.9",
		"jkcfg/kubernetes:0.6.2",
		"localhost:5000/foo:master",
		"jkcfg/kubernetes@sha256:93086646d44cd87edc3cc5b9c48133f36db122683bf5865d3e353e3d132a85bd",
	}

	for _, arg := range positives {
		t.Run(arg, func(t *testing.T) {
			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			var refs []name.Reference
			val := NewImageRefSliceValue(&refs)
			// it's just easier to get a flagset to create flag for us
			flags.Var(val, "lib", "a test flag")
			assert.NoError(t, flags.Parse([]string{"--lib", arg}))
			assert.Equal(t, val.getStringSlice(), []string{arg})
		})
	}

	t.Run("slice parsing", func(t *testing.T) {
		// Check that more than one can be passed
		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
		var refs []name.Reference
		val := NewImageRefSliceValue(&refs)
		// it's just easier to get a flagset to create flag for us
		flags.Var(val, "lib", "a test flag")
		args := []string{
			"--lib", positives[0],
			"--lib", positives[1] + "," + positives[2],
		}
		assert.NoError(t, flags.Parse(args))
		assert.Equal(t, val.getStringSlice(), positives[:3])
	})

	negatives := []string{
		"foo",                     // no tag
		"foo/bar:",                // still no tag; malformed
		"punctation/'notallowed'", // illegal characters
		"foo@sha256:abc123",       // malformed digest
	}

	for _, arg := range negatives {
		t.Run("[bad]"+arg, func(t *testing.T) {
			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			var refs []name.Reference
			// it's just easier to get a flagset to create flag for us
			flags.Var(NewImageRefSliceValue(&refs), "lib", "a test flag")
			assert.Error(t, flags.Parse([]string{"--lib", arg}))
		})
	}

	// Check that a bad apple results in a parse failure
	t.Run("bad apple", func(t *testing.T) {
		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
		var refs []name.Reference
		val := NewImageRefSliceValue(&refs)
		// it's just easier to get a flagset to create flag for us
		flags.Var(val, "lib", "a test flag")
		args := []string{
			"--lib", positives[1],
			"--lib", positives[2],
			"--lib", negatives[0],
		}
		assert.Error(t, flags.Parse(args))
	})
}
