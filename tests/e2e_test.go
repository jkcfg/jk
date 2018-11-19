package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func find(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, strings.TrimPrefix(path, dir))
		return nil
	})

	return files, err
}

func basename(testFile string) string {
	ext := filepath.Ext(testFile)
	return testFile[:len(testFile)-len(ext)]
}

func runTest(t *testing.T, file string) {
	base := basename(file)
	expectedDir := base + ".expected"
	gotDir := base + ".got"

	cmd := exec.Command("jk", "run", "-o", gotDir, file)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	// 1. Compare stdout/err.
	expected, _ := ioutil.ReadFile(file + ".expected")
	assert.Equal(t, string(expected), string(output))

	// 2. Compare produced files.
	expectedFiles, _ := find(expectedDir)
	gotFiles, _ := find(gotDir)

	// 2. a) Compare the list of files.
	assert.Equal(t, expectedFiles, gotFiles)

	// 2. b) Compare file content.
	for i := range expectedFiles {
		expected, err := ioutil.ReadFile(expectedDir + expectedFiles[i])
		assert.NoError(t, err)
		got, err := ioutil.ReadFile(gotDir + gotFiles[i])
		assert.NoError(t, err)

		assert.Equal(t, string(expected), string(got))
	}
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
