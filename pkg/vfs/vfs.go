// Package vfs contains abstractions for working with virtual
// filesystems. These are used to abstract access to files that may
// come from several different places, and to limit the ability to
// read files outside of those places expressly allowed.
package vfs

import (
	"net/http"
	"os"
)

// Location is a path within a specific filesystem.
type Location struct {
	Vfs  http.FileSystem
	Path string
}

// Nowhere is a zero value to use when e.g,. module resolution fails.
var Nowhere = Location{}

// EmptyFileSystem is an http.FileSystem with no files.
type EmptyFileSystem struct{}

// Open implements http.FileSystem
func (empty EmptyFileSystem) Open(p string) (http.File, error) {
	return nil, &os.PathError{Op: "open", Path: p, Err: os.ErrNotExist}
}

// Empty is an empty filesystem
var Empty = EmptyFileSystem{}
