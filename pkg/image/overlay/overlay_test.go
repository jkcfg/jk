package overlay

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func newMapFS(files map[string]string) http.FileSystem {
	return httpfs.New(mapfs.New(files))
}

// Test the simplest case: can I read a file from a single layer
func TestSingleLayerSingleFile(t *testing.T) {
	files := map[string]string{"dir/foo": "here is the text"}
	fs := New(newMapFS(files))
	assert.NotNil(t, fs)

	for _, p := range []string{"dir/foo", "./dir/foo", "/dir/foo", "/./dir/foo"} {
		file, err := fs.Open(p)
		assert.NoError(t, err)
		assert.NotNil(t, file)
		bytes, err := ioutil.ReadAll(file)
		assert.NoError(t, err)
		assert.Equal(t, files["dir/foo"], string(bytes))
	}
}

// Test the second simplest case: can I read a file that is not in the
// uppermost layer.
func TestSingleFile(t *testing.T) {
	files1 := map[string]string{"bar": "not the text"}
	files2 := map[string]string{"foo": "here is the text"}

	fs := New(newMapFS(files1), newMapFS(files2))
	assert.NotNil(t, fs)
	file, err := fs.Open("foo")
	assert.NoError(t, err)
	assert.NotNil(t, file)
	bytes, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, files2["foo"], string(bytes))
}

// Mainly a check that I know what I'm doing when using mapfs
func TestMapReaddir(t *testing.T) {
	files1 := map[string]string{
		"dir/foo": "unimportant",
		"dir/bar": "unimportant",
	}
	fs := newMapFS(files1)
	assert.NotNil(t, fs)

	dir, err := fs.Open("/dir")
	assert.NoError(t, err)

	info, err := dir.Stat()
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	_, err = dir.Read(nil)
	assert.Error(t, err)

	infos, err := dir.Readdir(0)
	assert.NoError(t, err)
	assert.Len(t, infos, 2)
}

// Test that readdir will return all the appropriate entries, that is
// directory entries across all layers, including one that does not
// contain the directory, and one that defines the file in question as
// a non-directory (stopping the search).
func TestReaddir(t *testing.T) {
	files1 := map[string]string{
		"dir/foo": "unimportant",
		"dir/bar": "unimportant",
	}
	// a layer in which no such dir exists
	files2 := map[string]string{
		"other/foo": "unimportant",
	}
	files3 := map[string]string{
		"dir/baz": "unimportant",
		"dir/bop": "unimportant",
	}
	files4 := map[string]string{
		"dir": "dead end",
	}
	files5 := map[string]string{
		"dir/boo": "unimportant",
	}
	fs := New(
		newMapFS(files1),
		newMapFS(files2),
		newMapFS(files3),
		newMapFS(files4),
		newMapFS(files5),
	)
	assert.NotNil(t, fs)

	dir, err := fs.Open("/dir")
	assert.NoError(t, err)

	info, err := dir.Stat()
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	_, err = dir.Read(nil)
	assert.Error(t, err)

	infos, err := dir.Readdir(0)
	assert.NoError(t, err)
	assert.Len(t, infos, 4)
}

// Test that readdir will discard duplicate entries; these would not
// correspond to files that can be opened, since the file is present
// on an upper layer.
func TestReadDuplicates(t *testing.T) {
	files1 := map[string]string{
		"dir/foo": "unimportant",
		"dir/bar": "unimportant",
	}
	files2 := map[string]string{
		"dir/foo": "overridden by upper layer",
		"dir/baz": "fresh file",
	}

	fs := New(newMapFS(files1), newMapFS(files2))
	assert.NotNil(t, fs)

	dir, err := fs.Open("/dir")
	assert.NoError(t, err)

	info, err := dir.Stat()
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	_, err = dir.Read(nil)
	assert.Error(t, err)

	infos, err := dir.Readdir(0)
	assert.NoError(t, err)
	assert.Len(t, infos, 3)
}

// Test that whiteout files can be used to "delete" files in a lower
// layer.
func TestWhiteout(t *testing.T) {
	files1 := map[string]string{
		"dir/.wh.foo":        "empty", // hides the file in the lower layer
		"dir/.wh.bar":        "empty", // hides file in lower layer, but not this layer
		"dir/bar":            "not hidden",
		"other/.wh.dir":      "empty", // hides a directory in the lower layer
		"third/.wh..wh..opq": "empty", // hides everything under third/ in lower layers
	}
	files2 := map[string]string{
		"dir/foo":                 "hidden",
		"other/dir/bar":           "hidden",
		"third/deep/path/to/file": "hidden",
	}
	fs := New(newMapFS(files1), newMapFS(files2))
	assert.NotNil(t, fs)

	// Check the directory exists
	dir, err := fs.Open("/dir")
	assert.NoError(t, err)
	info, err := dir.Stat()
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	// Check that the whiteout-deleted file can't be seen
	_, err = fs.Open("/dir/foo")
	assert.Error(t, err)
	// assert.True(t, errors.Is(err, os.ErrNotExist)) // TODO go 1.13

	// Check that the whiteout file can't be seen
	_, err = fs.Open("/dir/.wh.foo")
	assert.Error(t, err)
	// assert.True(t, errors.Is(err, os.ErrNotExist)) // TODO go 1.13

	// Check that the file with the whiteout _in the same layer_ can
	// be seen
	_, err = fs.Open("/dir/bar")
	assert.NoError(t, err)

	// Check that directory listing works, but with no results
	_, err = dir.Read(nil)
	assert.Error(t, err)
	infos, err := dir.Readdir(0)
	assert.NoError(t, err)
	assert.Len(t, infos, 1)

	// Opaque whiteout: check the the directory exists but is empty
	dir, err = fs.Open("/third")
	assert.NoError(t, err)
	info, err = dir.Stat()
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
	infos, err = dir.Readdir(0)
	assert.NoError(t, err)
	assert.Len(t, infos, 0)

	// and that the opaque whiteout file can't be seen
	_, err = fs.Open("/third/.wh..wh..opq")
	assert.Error(t, err)
	// assert.True(t, errors.Is(err, os.ErrNotExist)) // TODO go 1.13

	// and that we can't directly open the file
	_, err = fs.Open("/third/deep/path/to/file")
	assert.Error(t, err)
	// assert.True(t, errors.Is(err, os.ErrNotExist)) // TODO go 1.13
}
