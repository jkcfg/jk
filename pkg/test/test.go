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

func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func find(dir string) ([]string, error) {
	var files []string

	if !isDir(dir) {
		return nil, nil
	}

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

// Namer names things. Namer functions take the script under test as intput and
// returns the name. Namer functions are used to customize how files and
// directories are named.
type Namer func(script string) string

// Options are options that can be specified when creating a Test.
type Options struct {
	// Name is the test name. This value is used as the go test name. When left
	// empty, the script file name is used to derive the test name.
	Name string

	// WorkingDirectory is a directory to change to before executing the test.
	WorkingDirectory string

	// Env is additional environment entries for the command run.
	// These will be expanded in the commands themselves, if they are
	// supplied in a *.cmd file, since it is executed via `sh`.
	Env []string

	// CommandFile is a Namer that returns how files containing the list of
	// commands to run for a test should be named. It defaults to $script.cmd.
	CommandFile Namer

	// OutputDirectory is a Namer that returns how to name the directory the script
	// should output generated files to. It defaults to:
	//   echo $(echo $script | cut -f 1 -d '.').got
	OutputDirectory Namer

	// ExpectedOutputFile is a Namer that returns how files containing the expected
	// stdout should be named. It defaults to $script.expected.
	ExpectedOutputFile Namer

	// ExpectedOutputDirectory is a Namer that returns how the directory container
	// the expected generated names should be named. It defaults to:
	//   echo $(echo $script | cut -f 1 -d '.').expected
	ExpectedOutputDirectory Namer
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

func defaultOutputDirectory(script string) string {
	return basename(script) + ".got"
}

func (test *Test) outputDirectory() string {
	namer := test.opts.OutputDirectory
	if namer == nil {
		namer = defaultOutputDirectory
	}
	return namer(test.file)
}

func defaultCommandFile(script string) string {
	return script + ".cmd"
}

func (test *Test) commandFile() string {
	namer := test.opts.CommandFile
	if namer == nil {
		namer = defaultCommandFile
	}
	return namer(test.file)
}

func defaultExpectedOutputFile(script string) string {
	return script + ".expected"
}

func (test *Test) expectedOutputFile() string {
	namer := test.opts.ExpectedOutputFile
	if namer == nil {
		namer = defaultExpectedOutputFile
	}
	return namer(test.file)
}

func defaultExpectedOutputDirectory(script string) string {
	return basename(script) + ".expected"
}

func (test *Test) expectedOutputDirectory() string {
	namer := test.opts.ExpectedOutputDirectory
	if namer == nil {
		namer = defaultExpectedOutputDirectory
	}
	return namer(test.file)
}

func (test *Test) env() []string {
	if test.opts.Env != nil {
		return append(os.Environ(), test.opts.Env...)
	}
	return nil
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
		"%d", test.outputDirectory(),
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

	f, err := os.Open(test.commandFile())
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		args := test.parseCmd(scanner.Text())
		cmd := exec.Command("/bin/sh", "-c", strings.Join(args, " "))
		cmd.Dir = test.opts.WorkingDirectory
		cmd.Env = test.env()
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
	cmd := exec.Command("jk", "run", "-o", test.outputDirectory(), test.file)
	cmd.Dir = test.opts.WorkingDirectory
	cmd.Env = test.env()
	if err := test.setStdin(cmd); err != nil {
		return "", err
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (test *Test) exec() (string, error) {
	if exists(test.commandFile()) {
		return test.execWithCmd()
	}
	return test.execDefault()
}

// Run executes the test and compare its output to the expected state.
func (test *Test) Run(t *testing.T) {
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
	expected, _ := ioutil.ReadFile(test.expectedOutputFile())
	assert.Equal(t, string(expected), string(output))

	// 2. Compare produced files.
	expectedFiles, _ := find(test.expectedOutputDirectory())
	gotFiles, _ := find(test.outputDirectory())

	// 2. a) Compare the list of files.
	if !assert.Equal(t, expectedFiles, gotFiles) {
		assert.FailNow(t, "generated files not equivalent; bail")
	}

	// 2. b) Compare file content.
	for i := range expectedFiles {
		expected, err := ioutil.ReadFile(test.expectedOutputDirectory() + expectedFiles[i])
		assert.NoError(t, err)
		got, err := ioutil.ReadFile(test.outputDirectory() + gotFiles[i])
		assert.NoError(t, err)

		assert.Equal(t, string(expected), string(got))
	}
}
