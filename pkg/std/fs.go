package std

import (
	"errors"
	"path"
	"sort"

	"github.com/jkcfg/jk/pkg/vfs"

	"github.com/shurcooL/httpfs/vfsutil"
)

// FileInfo is the result from a std.fileinfo RPC (and used to
// represent each file, within Directory)
type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isdir"`
}

// Directory is the result from an std.dir RPC
type Directory struct {
	Name  string     `json:"name"`
	Path  string     `json:"path"`
	Files []FileInfo `json:"files"`
}

// MakeFileInfo returns a response to a FileInfo request, encoded
// ready to send to the V8 worker.
func MakeFileInfo(r Sandbox, path, module string) (FileInfo, error) {
	loc, err := r.getReadPath(path, module)
	if err != nil {
		return FileInfo{}, err
	}
	return fileInfo(loc, path)
}

// MakeDirectoryListing returns a response to a Dir request, encoded
// ready to send to the V8 worker.
func MakeDirectoryListing(r Sandbox, p, module string) (Directory, error) {
	loc, err := r.getReadPath(p, module)
	if err != nil {
		return Directory{}, err
	}
	return directoryListing(loc, p)
}

func fileInfo(loc vfs.Location, p string) (FileInfo, error) {
	info, err := vfsutil.Stat(loc.Vfs, loc.Path)
	switch {
	case err != nil:
		return FileInfo{}, err
	case !(info.IsDir() || info.Mode().IsRegular()):
		return FileInfo{}, errors.New("not a regular file")
	}
	return FileInfo{Name: info.Name(), Path: p, IsDir: info.IsDir()}, nil
}

func directoryListing(loc vfs.Location, p string) (Directory, error) {
	dir, err := loc.Vfs.Open(loc.Path)
	if err != nil {
		return Directory{}, err
	}
	infos, err := dir.Readdir(0)
	if err != nil {
		return Directory{}, err
	}

	// Sort the fileinfos by name, to avoid introducing non-determinism
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name() < infos[j].Name()
	})

	files := make([]FileInfo, 0, len(infos))

	for i := range infos {
		if infos[i].IsDir() || infos[i].Mode().IsRegular() {
			files = append(files, FileInfo{
				Name:  infos[i].Name(),
				Path:  path.Join(p, infos[i].Name()),
				IsDir: infos[i].IsDir(),
			})
		}
	}

	return Directory{
		Name:  path.Base(p),
		Path:  p,
		Files: files,
	}, nil
}
