package tests

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runTest(t *testing.T, file string) {
	expected, err := ioutil.ReadFile(file + ".expected")
	assert.NoError(t, err)

	cmd := exec.Command("jk", "run", file)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(output))
}

func TestEndToEnd(t *testing.T) {
	files, err := filepath.Glob("test-*.js")
	assert.NoError(t, err)

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			runTest(t, file)
		})
	}
}
