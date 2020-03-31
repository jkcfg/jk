// Package cache has code for downloading and caching container
// images, in such a way that each can be used as overlay filesystems.
//
// Container images come as a set of layers as gzipped tarballs, and a
// manifest, which lists the layers by name. The name of a layer is
// its (SHA256) digest.
//
// Since layers are content-addressed and can be shared between
// images, the layers from all images are lumped together into a
// directory.
//
// Manifests are stored in a file named for its digest; tags are
// symlinked to the "real" file.
//
//
package cache

// Sketch of how we get from a `--lib image:tag` argument to an overlay
// filesystem:
//
//   1. canonicalise the image name, and look in the filesystem for the
//   symlink or file named for the digest, in images/
//     - if not present, resolve the image name/tag/digest and get the
//     manifest from the registry
//   2. parse the manifest, and for each layer,
//      2.1. look for the layer in the layers/ directory
//           - if not present, download it and verify it
//   3. construct an overlay filesytem out of the layers

// Ref for downloading things:
// https://github.com/containers/skopeo/blob/master/cmd/skopeo/copy.go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	oci_v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/jkcfg/jk/pkg/image/overlay"
	"github.com/jkcfg/jk/pkg/vfs"
)

const (
	appDir       = "jk"
	layersDir    = "layers"
	manifestsDir = "manifests"
)

// Cache is a base directory for an image cache
type Cache struct {
	base string
}

// New constructs an instance of Cache given a cache
// directory. Usually the cache directory would be obtained with
// `os.UserCacheDir()`.
func New(userCacheDir string) *Cache {
	return &Cache{base: filepath.Join(userCacheDir, appDir)}
}

// Standardises the construction of the path to a layer.
func (cache *Cache) layerPath(algo, digest string) string {
	return filepath.Join(cache.base, layersDir, algo, digest)
}

// Standardises the construction of the path to a manifest.
func (cache *Cache) manifestPath(image, ref string) string {
	return filepath.Join(cache.base, manifestsDir, image, ref)
}

// FileSystemForImage takes an image name and ref (tag), and
// constructs a vfs.FileSystem from the image's layers as found in the
// cache. It assumes the manifest and layers will be present in the
// cache.
func (cache *Cache) FileSystemForImage(image, ref string) (vfs.FileSystem, error) {
	m := cache.manifestPath(image, ref)
	mfile, err := os.Open(m)
	if err != nil {
		return nil, fmt.Errorf("cannot stat manifest at implied path %s: %s", m, err.Error())
	}

	// This is the manifest type:
	// https://github.com/opencontainers/image-spec/blob/master/specs-go/v1/manifest.go
	var manifest oci_v1.Manifest
	err = json.NewDecoder(mfile).Decode(&manifest)
	if err != nil {
		return nil, fmt.Errorf("cannot decode manifest at %s: %s", m, err.Error())
	}

	layerCount := len(manifest.Layers)
	layers := make([]http.FileSystem, layerCount, layerCount)
	for i, desc := range manifest.Layers {
		layerPath := cache.layerPath(desc.Digest.Algorithm().String(), desc.Digest.Encoded())
		info, err := os.Stat(layerPath)
		if err != nil {
			return nil, fmt.Errorf("layer is not in image cache %s", layerPath)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("layer is not a directory as expected %s", layerPath)
		}
		// the overlay filesystem gets its layers top-most first, but
		// the layers are bottom-most first in the manifest.
		layers[layerCount-i-1] = http.Dir(layerPath)
	}
	return vfs.User(image+":"+ref+"!", overlay.New(layers...)), nil
}
