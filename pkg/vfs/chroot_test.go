package vfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func TestChroot(t *testing.T) {
	files := map[string]string{
		"foo":      "here is a file outside the new root dir",
		"root/bar": "here is a file inside the new root dir",
	}
	fs := User("mem", httpfs.New(mapfs.New(files)))
	chfs := Chroot(fs, "/root/")

	t.Run("/root/bar is accessible as /bar", func(t *testing.T) {
		_, err := chfs.Open("/bar")
		assert.NoError(t, err)
	})

	t.Run("/root/bar is not accessible as /root/bar", func(t *testing.T) {
		_, err := chfs.Open("/root/bar")
		assert.Error(t, err)
	})

	t.Run("/foo is not accessible via a parent path ../foo", func(t *testing.T) {
		_, err := chfs.Open("../foo")
		assert.Error(t, err)
	})
}
