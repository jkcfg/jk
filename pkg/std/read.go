package std

import (
	//	"context"
	"io/ioutil"
)

func read(url string) ([]byte, error) {
	// TODO(michael): optionally (by default) check that the file is "here or down"
	return ioutil.ReadFile(url)
}
