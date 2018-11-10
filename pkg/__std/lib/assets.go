package lib

import (
	"io/ioutil"
)

// ReadAll reads all content from path.
func ReadAll(path string) ([]byte, error) {
	f, err := Assets.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
