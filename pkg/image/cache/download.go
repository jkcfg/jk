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
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// FIXME don't force permissions on the actual files, since some may
// need to be executable
const cacheFileMode = os.FileMode(0400)
const cacheDirMode = os.FileMode(0700)

// Download makes sure the manifest and layers for a particular image
// are present in the cache.
func (c *Cache) Download(image string) error { // <-- could return the digest?
	// Ref for this code:
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
			// TODO(michael): factor this out
			if err = os.MkdirAll(filepath.Dir(tagManifestPath), cacheDirMode); err != nil {
				return err
			}
			err = os.Symlink(manifestPath, tagManifestPath)
			return err
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
		// Check for the layer in the layers directory
		digest, err := layer.Digest()
		if err != nil {
			return err
		}
		p := c.layerPath(digest.Algorithm, digest.Hex)
		_, err = os.Stat(p)
		if err == nil { // already have it
			continue
		}

		// TODO(michael): if there's a problem after this, clean up
		// the expanded layer before returning, so it doesn't look
		// like we've succeeded. In fact, better to expand somewhere
		// else, then rename it to the right place if it succeeds.
		if err = os.MkdirAll(p, cacheDirMode); err != nil {
			return err
		}

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
			targetPath := filepath.Join(p, hdr.Name)
			info := hdr.FileInfo()
			if info.IsDir() {
				if os.MkdirAll(targetPath, 0755); err != nil {
					// TODO cleanup
					return err
				}
			} else {
				f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, cacheFileMode)
				if err != nil {
					return err // TODO cleanup
				}
				if _, err := io.Copy(f, rdr); err != nil {
					f.Close()
					return err // TODO cleanup
				}
				f.Close()
			}
		}
	}

	// TODO(michael) figure out if we will always get an OCI v1
	// (compatible) manifest
	man, err := img.RawManifest()
	if err != nil {
		// TODO cleanup
		return err
	}

	if err = os.MkdirAll(filepath.Dir(manifestPath), cacheDirMode); err != nil {
		return err
	}
	if err = ioutil.WriteFile(manifestPath, man, cacheFileMode); err != nil {
		// TODO cleanup
		return err
	}

	if tagManifestPath != "" {
		if err = os.MkdirAll(filepath.Dir(tagManifestPath), cacheDirMode); err != nil {
			return err
		}
		err = os.Symlink(manifestPath, tagManifestPath)
		return err
	}

	return nil
}
