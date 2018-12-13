package std

import (
	"io/ioutil"

	"golang.org/x/text/encoding/unicode"

	"github.com/jkcfg/jk/pkg/__std"
)

func read(url string, encoding __std.Encoding) ([]byte, error) {
	// TODO(michael): optionally (by default) check that the file is "here or down"
	bytes, err := ioutil.ReadFile(url)
	switch encoding {
	case __std.EncodingBytes:
		break
	case __std.EncodingUTF16:
		encoder := unicode.UTF16(NativeEndian, unicode.IgnoreBOM).NewEncoder()
		bytes, err = encoder.Bytes(bytes)
	}
	return bytes, err
}
