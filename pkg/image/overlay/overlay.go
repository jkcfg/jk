package overlay

// Treat a set of layers as a filesystem, for e.g., putting in the
// module search path.

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
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

const whiteoutPrefix = ".wh."
const whiteoutOpaque = ".wh..wh..opq"

func whiteoutPath(d, base string) string {
	return path.Join(d, whiteoutPrefix+base)
}

func isWhiteout(base string) bool {
	return strings.HasPrefix(base, whiteoutPrefix)
}

// Open implements http.FileSystem#Open for the overlay filesystem
func (o *Overlay) Open(p string) (http.File, error) {
	p = path.Clean("/" + p) // this to ensure we don't get e.g., /./ when calling path.Split later

	// These are bookkeeping for whiteout handling, and will be
	// destructively updated
	woDir, woBase := path.Split(p)

	// whiteout files are never visible
	if isWhiteout(woBase) {
		return nil, &os.PathError{Op: "open", Path: p, Err: os.ErrNotExist}
	}

	// Start off with the potential whiteout files in the current
	// directory; we may add more later
	whiteouts := []string{whiteoutPath(woDir, woBase), path.Join(woDir, whiteoutOpaque)}

	// hasWhiteout returns true if there's a whiteout file in _this_
	// layer that would hide the path in _lower_ layers (false
	// otherwise). There's some bookkeeping here, so that we don't
	// repeat path calculations.
	hasWhiteout := func(fs http.FileSystem) bool {
		// First, those paths we've already calculated
		for _, wo := range whiteouts {
			if _, err := fs.Open(wo); err == nil {
				return true
			}
		}
		// Second, in case we haven't already, any paths up to the
		// root (saving them for next time)
		for ; woDir != "/"; woDir, woBase = path.Split(woDir) {
			woDir = woDir[:len(woDir)-1] // drop trailing slash
			wo := path.Join(woDir, whiteoutOpaque)
			if _, err := fs.Open(wo); err == nil {
				return true
			}
			whiteouts = append(whiteouts, wo)
			wo = path.Join(woDir, whiteoutPrefix+woBase)
			if _, err := fs.Open(wo); err == nil {
				return true
			}
			whiteouts = append(whiteouts, wo)
		}
		return false
	}

	for i, fs := range o.layers {
		file, err := fs.Open(p)
		if err != nil {
			if hasWhiteout(fs) {
				break
			} else {
				continue
			}
		}

		stat, err := file.Stat()
		if err != nil {
			if hasWhiteout(fs) {
				break
			} else {
				continue
			}
		}

		// If it's not a directory, this hides any lower entries
		if !stat.IsDir() {
			return file, nil
		}

		// It _is_ a directory. This will have the entries of all
		// lower directory entries.
		return &dir{
			path:  p,
			File:  file,
			upper: o.layers[i],
			lower: o.layers[i+1:],
		}, nil
	}

	return nil, &os.PathError{Op: "open", Path: p, Err: os.ErrNotExist}
}

// dir represents a directory in an overlay filesystem
type dir struct {
	http.File // the uppermost file
	path      string
	upper     http.FileSystem
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
		whiteout := map[string]struct{}{}
		record := func(infos []os.FileInfo) {
			for _, info := range infos {
				name := info.Name()
				if isWhiteout(name) {
					whiteout[name[len(whiteoutPrefix):]] = struct{}{}
					continue
				}
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

		// calculate all the possible whiteout files that would hide
		// the directory in lower layers
		dir, base := path.Clean("/"+d.path), ""
		whiteouts := []string{path.Join(dir, whiteoutOpaque)}
		for ; dir != "/"; dir, base = path.Split(dir) {
			whiteouts = append(whiteouts, path.Join(dir, whiteoutPrefix+base), path.Join(dir, whiteoutOpaque))
		}

		current := d.upper
		for _, layer := range d.lower {
			// if any of these exist, we're done
			for _, wo := range whiteouts {
				if _, err := current.Open(wo); err == nil {
					return appropriate(nil)
				}
			}
			current = layer

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

			// make sure we ignore anything that was whited-out in a
			// prior layer
			for n, f := range whiteout {
				seen[n] = f
			}
			whiteout = map[string]struct{}{}
			infos, err := file.Readdir(0)
			record(infos)
			if err != nil {
				return appropriate(err)
			}
		}
	}

	return appropriate(nil)
}
