package std

import (
//	"context"
)

func read(url string) ([]byte, error) {
	return []byte("echoing " + url), nil
}
