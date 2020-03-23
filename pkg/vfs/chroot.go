package vfs

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type chroot struct {
	FileSystem
	newRoot string
}

// Open opens a file in the chroot-ed filesystem.
func (fs *chroot) Open(p string) (http.File, error) {
	p = path.Clean(p)
	if strings.HasPrefix(p, "../") {
		return nil, os.ErrNotExist
	}
	return fs.FileSystem.Open(path.Join(fs.newRoot, p))
}

func (fs *chroot) QualifyPath(p string) string {
	return fs.FileSystem.QualifyPath(path.Join(fs.newRoot, p))
}

// Chroot makes a given directory appear to be the root of the
// filesystem.
func Chroot(vfs FileSystem, newRoot string) FileSystem {
	return &chroot{
		FileSystem: vfs,
		newRoot:    newRoot,
	}
}
