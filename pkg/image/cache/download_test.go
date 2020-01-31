package cache

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func setupRegistry(t *testing.T) (*httptest.Server, string, string) {
	regHandler := registry.New()
	regSrv := httptest.NewServer(regHandler)
	img, err := crane.Load("./testfiles/helloworld.tar")
	assert.NoError(t, err)
	newImg := regSrv.URL[len("http://"):] + "/helloworld:linux"
	dig, err := img.Digest()
	assert.NoError(t, err)
	newImgDig := regSrv.URL[len("http://"):] + "/helloworld" + "@" + dig.String()
	assert.NoError(t, crane.Push(img, newImg))
	return regSrv, newImg, newImgDig
}

func TestDownloadToCache(t *testing.T) {
	tmp, err := ioutil.TempDir("", "jk-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)

	cache := New(tmp)

	regSrv, imgTag, imgDigest := setupRegistry(t)
	defer regSrv.Close()

	err = cache.Download(imgTag)
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
