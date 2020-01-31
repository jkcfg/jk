package cache

// The read part of the cache constructs a filesystem given an image
// ref. The download part, here, gets the constituents of images and
// puts them into the cache. The contract between the two is in the
// layout of the directories used by the cache, described in the
// package documentation.

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// FIXME don't force permissions on the actual files, since some may
// need to be executable
const cacheFileMode = os.FileMode(0400)
const cacheDirMode = os.FileMode(0700)

func linkTagManifest(digestManifestPath, tagManifestPath string) error {
	if err := os.MkdirAll(filepath.Dir(tagManifestPath), cacheDirMode); err != nil {
		return err
	}
	return os.Symlink(digestManifestPath, tagManifestPath)
}

// Download makes sure the manifest and layers for a particular image
// are present in the cache.
func (c *Cache) Download(image string) error {
	// Example of code using crane packages to pull images:
	// https://github.com/google/go-containerregistry/blob/master/pkg/crane/pull.go#Save
	ref, err := name.ParseReference(image)
	if err != nil {
		return err
	}

	manifestPath := c.manifestPath(ref)
	if manifestPath == "" {
		return fmt.Errorf("cannot make manifest path for image ref %q", ref)
	}

	// Whichever kind of image ref, if the manifest is already in the
	// filesystem, we must have completed this previously.
	_, err = os.Stat(manifestPath)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	img, err := remote.Image(ref)
	if err != nil {
		return err
	}

	// If it's a reference with a tag, specifically, then look for the
	// manifest at the digest as well
	var tagManifestPath string
	if _, ok := ref.(name.Tag); ok {
		dig, err := img.Digest()
		if err != nil {
			return err
		}
		digestManifestPath := c.manifestPath(ref.Context().Digest(dig.String()))
		if digestManifestPath == "" {
			return fmt.Errorf("unable to construct path to manifest for %q", dig)
		}
		// we'll be writing to the digest path and symlinking the tag
		// path to the digest path
		tagManifestPath, manifestPath = manifestPath, digestManifestPath
		_, err = os.Stat(manifestPath)
		if err == nil {
			return linkTagManifest(manifestPath, tagManifestPath)
		}
		if !os.IsNotExist(err) {
			return err
		}
	}

	// TODO any kind of verification

	layers, err := img.Layers()
	if err != nil {
		return err
	}
	for _, layer := range layers {
		if err = c.writeLayer(layer); err != nil {
			return err
		}
	}

	// TODO(michael) figure out if we will always get an OCI v1
	// (compatible) manifest
	man, err := img.RawManifest()
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(manifestPath), cacheDirMode); err != nil {
		return err
	}
	if err = ioutil.WriteFile(manifestPath, man, cacheFileMode); err != nil {
		return err
	}

	if tagManifestPath != "" {
		return linkTagManifest(manifestPath, tagManifestPath)
	}
	return nil
}

func (c *Cache) writeLayer(layer v1.Layer) error {
	// Check for the layer in the layers directory
	digest, err := layer.Digest()
	if err != nil {
		return err
	}
	layerPath := c.layerPath(digest.Algorithm, digest.Hex)
	_, err = os.Stat(layerPath)
	if err == nil { // already have it
		return nil
	}

	tmpLayerPath, err := ioutil.TempDir("", "jk.layer")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpLayerPath)

	// Write the layer by expanding it into the directory.
	layerReader, err := layer.Uncompressed()
	if err != nil {
		return err
	}
	rdr := tar.NewReader(layerReader)
	for {
		hdr, err := rdr.Next()
		if err == io.EOF {
			break
		}
		targetPath := filepath.Join(tmpLayerPath, hdr.Name)
		info := hdr.FileInfo()
		if info.IsDir() {
			if os.MkdirAll(targetPath, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		} else {
			f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, hdr.FileInfo().Mode())
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, rdr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	if err = os.MkdirAll(filepath.Dir(layerPath), cacheDirMode); err != nil {
		return err
	}
	if err = os.Rename(tmpLayerPath, layerPath); err != nil {
		return err
	}

	return nil
}
