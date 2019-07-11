package test

import (
	"bufio"
	"fmt"
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

// Options are options that can be specified when creating a Test.
type Options struct {
	// Name is the test name. This value is used as the go test name. When left
	// empty, the script file name is used to derive the test name.
	Name string

	// WorkingDirectory is a directoring to change to before executing the test.
	WorkingDirectory string
}

// Test is a end to end test, corresponding to one test-$testname.js file.
type Test struct {
	file string // name of the test-*.js test file
	opts Options
}

// New creates a new Test wrapping the given jk script.
func New(script string, options ...Options) *Test {
	test := &Test{
		file: script,
	}
	if len(options) > 0 {
		test.opts = options[0]
	}
	return test
}

func (test *Test) jsFile() string {
	return test.file
}

func (test *Test) basename() string {
	return basename(test.file)
}

// Name is the test name.
func (test *Test) Name() string {
	if test.opts.Name != "" {
		return test.opts.Name
	}
	if strings.HasPrefix(test.file, "test-") {
		return test.file[len("test-") : len(test.file)-3]
	}
	return test.file[:len(test.file)-3]
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (test *Test) shouldErrorOut() bool {
	return exists(test.file + ".error")
}

func (test *Test) shouldSkip() bool {
	return exists(test.file + ".skip")
}

func (test *Test) outputDir() string {
	return basename(test.file) + ".got"
}

func (test *Test) setStdin(cmd *exec.Cmd) error {
	if exists(test.file + ".in") {
		infile, err := os.Open(test.file + ".in")
		if err != nil {
			return err
		}
		cmd.Stdin = infile
	}
	return nil
}

func (test *Test) parseCmd(line string) []string {
	parts := strings.Split(line, " ")
	replacer := strings.NewReplacer(
		"%d", test.outputDir(),
		"%b", test.basename(),
		"%t", test.Name(),
		"%f", test.jsFile(),
	)
	// Replace special strings
	for i := range parts {
		parts[i] = replacer.Replace(parts[i])
	}

	return parts
}

func (test *Test) execWithCmd() (string, error) {
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
		cmd.Dir = test.opts.WorkingDirectory
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

func (test *Test) execDefault() (string, error) {
	cmd := exec.Command("jk", "run", "-o", test.outputDir(), test.file)
	cmd.Dir = test.opts.WorkingDirectory
	if err := test.setStdin(cmd); err != nil {
		return "", err
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (test *Test) exec() (string, error) {
	if exists(test.file + ".cmd") {
		return test.execWithCmd()
	}
	return test.execDefault()
}

// Run executes the test and compare its output to the expected state.
func (test *Test) Run(t *testing.T) {
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
