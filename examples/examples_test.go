package tests

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jkcfg/jk/pkg/test"
	"github.com/stretchr/testify/assert"
)

func listTestFiles(t *testing.T) []string {
	var files []string

	files, err := filepath.Glob("*/*/test*.cmd")
	assert.NoError(t, err)

	sort.Strings(files)
	return files
}

func wd(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join(parts[:2], "/")
}

func cmd(script string) string {
	return script
}

func expectedOutputFile(script string) string {
	basename := strings.TrimSuffix(filepath.Base(script), filepath.Ext(script))
	return filepath.Join(wd(script), basename+".expected")
}

func TestExamples(t *testing.T) {
	files := listTestFiles(t)

	for _, file := range files {
		test := test.New(file, test.Options{
			Name:               wd(file),
			WorkingDirectory:   wd(file),
			CommandFile:        cmd,
			ExpectedOutputFile: expectedOutputFile,
		})
		t.Run(test.Name(), func(t *testing.T) {
			test.Run(t)
		})
	}
}
