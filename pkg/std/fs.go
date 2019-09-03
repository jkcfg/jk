package std

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
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
	base, path, err := r.getPath(path, module)
	if err != nil {
		return FileInfo{}, err
	}
	return fileInfo(base, path)
}

// DirectoryListing returns a response to a Dir request, encoded ready
// to send to the V8 worker.
func (r ReadBase) DirectoryListing(path, module string) (Directory, error) {
	base, path, err := r.getPath(path, module)
	if err != nil {
		return Directory{}, err
	}
	return directoryListing(base, path)
}

func fileInfo(base, rel string) (FileInfo, error) {
	path := filepath.Join(base, rel)
	info, err := os.Stat(path)
	switch {
	case err != nil:
		return FileInfo{}, err
	case !(info.IsDir() || info.Mode().IsRegular()):
		return FileInfo{}, errors.New("not a regular file")
	}
	return FileInfo{Name: info.Name(), Path: rel, IsDir: info.IsDir()}, nil
}

func directoryListing(base, rel string) (Directory, error) {
	path := filepath.Join(base, rel)
	dir, err := os.Open(path)
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

	var files []FileInfo

	for i := range infos {
		if infos[i].IsDir() || infos[i].Mode().IsRegular() {
			files = append(files, FileInfo{
				Name:  infos[i].Name(),
				Path:  filepath.Join(rel, infos[i].Name()),
				IsDir: infos[i].IsDir(),
			})
		}
	}

	return Directory{
		Name:  filepath.Base(rel),
		Path:  rel,
		Files: files,
	}, nil
}
