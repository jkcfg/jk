// Package vfs contains abstractions for working with virtual
// filesystems. These are used to abstract access to files that may
// come from several different places, and to limit the ability to
// read files outside of those places expressly allowed.
package vfs

import (
	"net/http"
	"os"
	"path"
)

// FileSystem is an interface for filesystems used to source modules,
// files, and resources.
type FileSystem interface {
	http.FileSystem
	// IsInternal returns true if a filesystem should be considered
	// internal to the system, and therefore not recorded
	IsInternal() bool
	// QualifyPath gives a path within the filesystem a identifying ,
	// used both for reporting its location unambiguously (though not
	// necessarily using a valid path), and e.g., to avoid loading
	// modules more than once.
	QualifyPath(path string) string
}

type prefixed struct {
	prefix string
}

// QualifyPath takes a path assumed to be within the filesystem, and
// qualifies it using the prefix assigned.
func (f prefixed) QualifyPath(p string) string {
	return path.Join(f.prefix, p)
}

// UserFileSystem is a way to wrap a "regular" filesystem (e.g., as
// constructed by `http.Dir`) so it is marked as a non-system
// filesystem.
type UserFileSystem struct {
	prefixed
	http.FileSystem
}

// IsInternal implements the method of FileSystem, in the negative by
// definition.
func (u UserFileSystem) IsInternal() bool {
	return false
}

// User is a convenience for creating a user (non-system) filesystem
func User(prefix string, fs http.FileSystem) UserFileSystem {
	return UserFileSystem{
		prefixed:   prefixed{prefix},
		FileSystem: fs,
	}
}

// InternalFileSystem is a way to wrap a regular http.FileSystem so
// that it appears to be internal to the workings of the runtime,
// e.g., so reads don't get recorded.
type InternalFileSystem struct {
	prefixed
	http.FileSystem
}

// IsInternal implements the method of FileSystem, in the positive by
// definition.
func (f InternalFileSystem) IsInternal() bool {
	return true
}

// Internal is a convenience for creating an internal (system) filesystem
func Internal(prefix string, fs http.FileSystem) InternalFileSystem {
	return InternalFileSystem{
		prefixed:   prefixed{prefix},
		FileSystem: fs,
	}
}

// Location is a path within a specific filesystem.
type Location struct {
	Vfs  FileSystem
	Path string
}

// CanonicalPath gives an identifying (though not necessarily valid) path
// for the location, by including the filesystem's identity via
// IdentifyingPath.
func (loc Location) CanonicalPath() string {
	return loc.Vfs.QualifyPath(loc.Path)
}

// Nowhere is a zero value to use when e.g,. module resolution fails.
var Nowhere = Location{}

// EmptyFileSystem is an http.FileSystem with no files.
type EmptyFileSystem struct{}

// Open implements http.FileSystem
func (empty EmptyFileSystem) Open(p string) (http.File, error) {
	return nil, &os.PathError{Op: "open", Path: p, Err: os.ErrNotExist}
}

// Empty is an empty filesystem. It is marked as internal, since it is
// an implementation detail used to prevent relative resolution of
// (module) paths.
var Empty = Internal("<EMPTY>", EmptyFileSystem{})
