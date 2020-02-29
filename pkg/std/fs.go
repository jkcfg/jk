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

// FileInfo returns a response to a FileInfo request, encoded ready to
// send to the V8 worker.
func (r ReadBase) FileInfo(path, module string) (FileInfo, error) {
	loc, rel, err := r.getPath(path, module)
	if err != nil {
		return FileInfo{}, err
	}
	return fileInfo(loc, rel)
}

// DirectoryListing returns a response to a Dir request, encoded ready
// to send to the V8 worker.
func (r ReadBase) DirectoryListing(path, module string) (Directory, error) {
	loc, rel, err := r.getPath(path, module)
	if err != nil {
		return Directory{}, err
	}
	return directoryListing(loc, rel)
}

func fileInfo(loc vfs.Location, rel string) (FileInfo, error) {
	p := path.Join(loc.Path, rel)
	info, err := vfsutil.Stat(loc.Vfs, p)
	switch {
	case err != nil:
		return FileInfo{}, err
	case !(info.IsDir() || info.Mode().IsRegular()):
		return FileInfo{}, errors.New("not a regular file")
	}
	return FileInfo{Name: info.Name(), Path: rel, IsDir: info.IsDir()}, nil
}

func directoryListing(base vfs.Location, rel string) (Directory, error) {
	p := path.Join(base.Path, rel)
	dir, err := base.Vfs.Open(p)
	if err != nil {
		return Directory{}, err
	}
	defer dir.Close()
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
				Path:  path.Join(rel, infos[i].Name()),
				IsDir: infos[i].IsDir(),
			})
		}
	}

	return Directory{
		Name:  path.Base(rel),
		Path:  rel,
		Files: files,
	}, nil
}
