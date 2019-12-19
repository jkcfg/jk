package overlay

// Treat a set of layers as a filesystem, for e.g., putting in the
// module search path.

import (
	"io"
	"net/http"
	"os"
)

// Overlay is a read-only filesystem formed with a list of layers
// representing diffs. It is designed to behave like a filesystem such
// as those used in containers. The layers are given uppermost first
// (at index 0).
//
//     https://www.kernel.org/doc/Documentation/filesystems/overlayfs.txt
// is the inspiration for OCI image layers:
//
// > An overlay filesystem combines two filesystems - an 'upper' filesystem
// > and a 'lower' filesystem.  When a name exists in both filesystems, the
// > object in the 'upper' filesystem is visible while the object in the
// > 'lower' filesystem is either hidden or, in the case of directories,
// > merged with the 'upper' object.
type Overlay struct {
	layers []http.FileSystem
}

// New constructs an overlay filesystem given the layers
func New(layers ...http.FileSystem) *Overlay {
	return &Overlay{
		layers: layers,
	}
}

// Open implements http.FileSystem#Open for the overlay filesystem
func (o *Overlay) Open(path string) (http.File, error) {
	for i, fs := range o.layers {
		file, err := fs.Open(path)
		if err != nil {
			continue
		}

		stat, err := file.Stat()
		if err != nil {
			continue
		}

		// If it's not a directory, this hides any lower entries
		if !stat.IsDir() {
			return file, nil
		}

		// It _is_ a directory. This will have the entries of all
		// lower directory entries.
		return &dir{
			path:  path,
			File:  file,
			lower: o.layers[i+1:],
		}, nil
	}

	return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
}

// dir represents a directory in an overlay filesystem
type dir struct {
	http.File // the uppermost file
	path      string
	lower     []http.FileSystem
	// readdir state
	readdir []os.FileInfo
}

// Readdir implements http.File#Readdir for a directory in an overlay
// filesystem.
func (d *dir) Readdir(count int) ([]os.FileInfo, error) {
	// return the appropriate combination of slice and error, which
	// depends on the requested count. This assumes d.readdir has been
	// calculated.
	appropriate := func(err error) ([]os.FileInfo, error) {
		// This is a minor fudge: it may raise an error earlier than
		// it would otherwise been seen when using count > 0.
		if count > 0 {
			if err != nil {
				return nil, err
			}
			if len(d.readdir) == 0 {
				return nil, io.EOF
			}
			result := d.readdir[:count]
			d.readdir = d.readdir[count:]
			return result, nil
		}
		return d.readdir, err
	}

	// compile all the entries first
	if d.readdir == nil {
		seen := map[string]struct{}{}
		record := func(infos []os.FileInfo) {
			for _, info := range infos {
				name := info.Name()
				if _, ok := seen[name]; !ok {
					seen[name] = struct{}{}
					d.readdir = append(d.readdir, info)
				}
			}
		}

		infos, err := d.File.Readdir(0)
		record(infos)
		if err != nil {
			return appropriate(err)
		}

		for _, layer := range d.lower {
			file, err := layer.Open(d.path)
			if err != nil { // no file, check the next layer
				continue
			}
			info, err := file.Stat()
			if err != nil {
				return appropriate(err)
			}
			if !info.IsDir() {
				// a dead end, stop here.
				return appropriate(nil)
			}
			infos, err := file.Readdir(0)
			record(infos)
			if err != nil {
				return appropriate(err)
			}
		}
	}

	return appropriate(nil)
}
