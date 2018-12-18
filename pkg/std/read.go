package std

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/unicode"

	"github.com/jkcfg/jk/pkg/__std"

	"github.com/ghodss/yaml"
)

func readYAML(in io.Reader) ([]byte, error) {
	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	return yaml.YAMLToJSON(bytes)
}

func readJSON(in io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(in, &buf)
	var obj interface{}
	if err := json.NewDecoder(tee).Decode(&obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func readRaw(in io.Reader) ([]byte, error) {
	return ioutil.ReadAll(in)
}

type readFunc func(io.Reader) ([]byte, error)

func readerByPath(path string) readFunc {
	switch filepath.Ext(path) {
	case ".yml", ".yaml":
		return readYAML
	case ".json":
		return readJSON
	}
	return readJSON
}

func read(path string, format __std.Format, encoding __std.Encoding) ([]byte, error) {
	// TODO(michael): optionally (by default) check that the file is "here or down"
	var reader readFunc = readRaw

	if encoding == __std.EncodingJSON {
		switch format {
		case __std.FormatAuto:
			reader = readerByPath(path)
		case __std.FormatYAML:
			reader = readYAML
		case __std.FormatJSON:
			reader = readJSON
		default:
			reader = readJSON
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := reader(f)
	if err != nil {
		return nil, err
	}

	switch encoding {
	case __std.EncodingBytes:
		break
	case __std.EncodingString, __std.EncodingJSON:
		encoder := unicode.UTF16(NativeEndian, unicode.IgnoreBOM).NewEncoder()
		bytes, err = encoder.Bytes(bytes)
	}
	return bytes, err
}
