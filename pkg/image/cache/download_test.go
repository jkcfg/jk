package cache

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func setupRegistry(t *testing.T) *httptest.Server {
	regHandler := registry.New()
	regSrv := httptest.NewServer(regHandler)
	return regSrv
}

func fixture(t *testing.T, reg *httptest.Server, imageName string) (tag, digest string) {
	img, err := crane.Load("./testfiles/" + imageName + ".tar")
	assert.NoError(t, err)
	newImg := reg.URL[len("http://"):] + "/" + imageName + ":v1"
	dig, err := img.Digest()
	assert.NoError(t, err)
	newImgDig := reg.URL[len("http://"):] + "/" + imageName + "@" + dig.String()
	assert.NoError(t, crane.Push(img, newImg))
	return newImg, newImgDig
}

func TestDownloadToCache(t *testing.T) {
	tmp, err := ioutil.TempDir("", "jk-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)

	cache := New(tmp)

	regSrv := setupRegistry(t)
	defer regSrv.Close()

	imgTag, imgDigest := fixture(t, regSrv, "helloworld")

	err = cache.Download(mustParseRef(imgTag))
	assert.NoError(t, err)

	ov, err := cache.FileSystemForImage(mustParseRef(imgTag))
	assert.NoError(t, err)
	f, err := ov.Open("/hello")
	assert.NoError(t, err)
	defer f.Close()

	_, err = ioutil.ReadAll(f)
	assert.NoError(t, err)

	// The registry is switched off so we have to get it from the
	// cache.
	regSrv.Close()

	// Make sure we can get the image using its digest too.

	ov, err = cache.FileSystemForImage(mustParseRef(imgDigest))
	assert.NoError(t, err)
	f, err = ov.Open("/hello")
	assert.NoError(t, err)
	defer f.Close()

	_, err = ioutil.ReadAll(f)
	assert.NoError(t, err)
}

func TestDownloadWithSymlink(t *testing.T) {
	tmp, err := ioutil.TempDir("", "jk-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)

	cache := New(tmp)

	regSrv := setupRegistry(t)
	defer regSrv.Close()

	imgTag, _ := fixture(t, regSrv, "symlink")

	assert.NoError(t, cache.Download(mustParseRef(imgTag)))
	ov, err := cache.FileSystemForImage(mustParseRef(imgTag))
	assert.NoError(t, err)

	// The image has two files, symlink/v1/index.js and symlink/v2/index.js.
	// symlink/v2/index.js is a symlink to ../v1/index.js.
	// Test that:

	//  1. symlink/v2/index.js can be opened
	f, err := ov.Open("/symlink/v2/index.js")
	assert.NoError(t, err)
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, string(bytes), "export default 'module content';\n")

	//  2. the file symlink/v2/index.js in the layer is still a symlink
	man, err := cache.manifest(mustParseRef(imgTag))
	assert.NoError(t, err)
	// this is more to make sure the fixture image hasn't changed shape
	assert.Len(t, man.Layers, 1)

	desc := man.Layers[0]
	layerp := cache.layerPath(desc.Digest.Algorithm().String(), desc.Digest.Encoded())
	linkp := filepath.Join(layerp, "symlink", "v2", "index.js")
	_, err = os.Readlink(linkp) // this will error with "invalid argument" if it's not a link
	assert.NoError(t, err)
}
