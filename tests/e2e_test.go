package tests

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/assert"

	"github.com/jkcfg/jk/pkg/test"
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

// Run an image registry, loading any tar files that are in testfiles/
func runRegistry(t *testing.T) *httptest.Server {
	regHandler := registry.New()
	regSrv := httptest.NewServer(regHandler)
	host := regSrv.URL[len("http://"):]
	println("Registry server: ", regSrv.URL)
	files, err := filepath.Glob("testfiles/*.tar")
	assert.NoError(t, err)
	for _, file := range files {
		img, err := crane.Load(file)
		imageName := filepath.Base(file[:len(file)-len(".tar")]) + ":v1"
		assert.NoError(t, err)
		assert.NoError(t, crane.Push(img, host+"/"+imageName))
		println("Uploaded", imageName)
	}
	return regSrv
}

func TestEndToEnd(t *testing.T) {
	files := listTestFiles(t)

	reg := runRegistry(t)
	defer reg.Close()

	env := []string{
		"REGISTRY=" + reg.URL[len("http://"):],
	}

	tmp, err := ioutil.TempDir("", "jk-testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	for _, file := range files {
		testTmp := filepath.Join(tmp, file+".d")
		if err := os.Mkdir(testTmp, os.FileMode(0755)); err != nil {
			t.Fatal(err)
		}
		testEnv := append(env, "TEMP="+testTmp)
		test := test.New(file, test.Options{Env: testEnv})
		t.Run(test.Name(), func(t *testing.T) {
			test.Run(t)
		})
	}
}
