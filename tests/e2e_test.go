package tests

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/jkcfg/jk/pkg/test"
	"github.com/stretchr/testify/assert"
)

func listTestFiles(t *testing.T) []string {
	// Some tests aren't actually in this directory, but a .cmd file is used to
	// tune how jk is run. We need to account for those, making sure tests with
	// both a test-*.js file and a .cmd file aren't run twice.
	cmds, err := filepath.Glob("test-*.js.cmd")
	assert.NoError(t, err)

	files, err := filepath.Glob("test-*.js")
	assert.NoError(t, err)

	for _, cmd := range cmds {
		// Remove .cmd extension
		files = append(files, cmd[:len(cmd)-4])
	}

	// Deduplicate test files
	unique := make(map[string]struct{})
	for _, key := range files {
		unique[key] = struct{}{}
	}

	files = make([]string, 0, len(unique))
	for key := range unique {
		files = append(files, key)
	}

	sort.Strings(files)
	return files
}

func TestEndToEnd(t *testing.T) {
	files := listTestFiles(t)

	for _, file := range files {
		test := test.New(file)
		t.Run(test.Name(), func(t *testing.T) {
			test.Run(t)
		})
	}
}
