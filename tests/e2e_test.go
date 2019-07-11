package tests

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func find(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		case strings.HasSuffix(path, "~"):
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

// test is a end to end test, corresponding to one test-$testname.js file.
type test struct {
	file string // name of the test-*.js test file
}

func (test *test) jsFile() string {
	return test.file
}

func (test *test) basename() string {
	return basename(test.file)
}

func (test *test) name() string {
	return test.file[len("test-") : len(test.file)-3]
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (test *test) shouldErrorOut() bool {
	return exists(test.file + ".error")
}

func (test *test) shouldSkip() bool {
	return exists(test.file + ".skip")
}

func (test *test) outputDir() string {
	return basename(test.file) + ".got"
}

func (test *test) setStdin(cmd *exec.Cmd) error {
	if exists(test.file + ".in") {
		infile, err := os.Open(test.file + ".in")
		if err != nil {
			return err
		}
		cmd.Stdin = infile
	}
	return nil
}

func (test *test) parseCmd(line string) []string {
	parts := strings.Split(line, " ")
	replacer := strings.NewReplacer(
		"%d", test.outputDir(),
		"%b", test.basename(),
		"%t", test.name(),
		"%f", test.jsFile(),
	)
	// Replace special strings
	for i := range parts {
		parts[i] = replacer.Replace(parts[i])
	}

	return parts
}

func (test *test) execWithCmd() (string, error) {
	jkOutput := ""

	f, err := os.Open(test.file + ".cmd")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		args := test.parseCmd(scanner.Text())
		cmd := exec.Command("/bin/sh", "-c", strings.Join(args, " "))
		if err := test.setStdin(cmd); err != nil {
			return "", err
		}
		if args[0] == "jk" {
			output, err := cmd.CombinedOutput()
			if err != nil {
				// Display the output of jk in case of failure.
				fmt.Print(string(output))
				return "", err
			}
			jkOutput = string(output)
		} else {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return "", err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return jkOutput, nil
}

func (test *test) execDefault() (string, error) {
	cmd := exec.Command("jk", "run", "-o", test.outputDir(), test.file)
	if err := test.setStdin(cmd); err != nil {
		return "", err
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (test *test) exec() (string, error) {
	if exists(test.file + ".cmd") {
		return test.execWithCmd()
	}
	return test.execDefault()
}

func (test *test) run(t *testing.T) {
	base := basename(test.file)
	expectedDir := base + ".expected"
	gotDir := base + ".got"

	if test.shouldSkip() {
		return
	}

	output, err := test.exec()

	// 0. Check process exit code.
	if test.shouldErrorOut() {
		assert.Error(t, err)
		if err != nil {
			_, ok := err.(*exec.ExitError)
			assert.True(t, ok, err.Error())
		}
	} else {
		if err != nil {
			fmt.Print(string(output))
		}
		assert.NoError(t, err)
	}

	// 1. Compare stdout/err.
	expected, _ := ioutil.ReadFile(test.file + ".expected")
	assert.Equal(t, string(expected), string(output))

	// 2. Compare produced files.
	expectedFiles, _ := find(expectedDir)
	gotFiles, _ := find(gotDir)

	// 2. a) Compare the list of files.
	if !assert.Equal(t, expectedFiles, gotFiles) {
		assert.FailNow(t, "generated files not equivalent; bail")
	}

	// 2. b) Compare file content.
	for i := range expectedFiles {
		expected, err := ioutil.ReadFile(expectedDir + expectedFiles[i])
		assert.NoError(t, err)
		got, err := ioutil.ReadFile(gotDir + gotFiles[i])
		assert.NoError(t, err)

		assert.Equal(t, string(expected), string(got))
	}
}

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
		test := &test{file}
		t.Run(test.name(), func(t *testing.T) {
			test.run(t)
		})
	}
}
