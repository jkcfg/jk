package std

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dlespiau/jk/pkg/__std"

	"github.com/ghodss/yaml"
)

type closer func()

func nilCloser() {}

type writerFunc func(io.Writer, []byte)

func writeJSON(w io.Writer, value []byte) {
	var v interface{}
	if err := json.Unmarshal(value, &v); err != nil {
		log.Fatalf("writeJSON: unmarshal: %s", err)
	}
	i, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("writeJSON: marshal: %s", err)
	}
	w.Write(i)
	w.Write([]byte{'\n'})
}

func writeYAML(w io.Writer, value []byte) {
	y, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		log.Fatalf("writeYAML: %s", err)
	}
	w.Write(y)
}

func writer(path string) (io.Writer, closer) {
	if path == "" {
		return os.Stdout, nilCloser
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0666); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return f, func() { f.Close() }
}

func writerFuncFromPath(path string) writerFunc {
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml":
		fallthrough
	case ".yml":
		return writeYAML
	case ".json":
		fallthrough
	default:
		return writeJSON
	}
}

func write(value []byte, path string, format __std.OutputFormat) {
	w, close := writer(path)

	var out writerFunc
	switch format {
	case __std.OutputFormatAuto:
		out = writerFuncFromPath(path)
	case __std.OutputFormatJSON:
		out = writeJSON
	case __std.OutputFormatYAML:
		out = writeYAML
	default:
		log.Fatalf("write: unknown output format (%d)", int(format))
	}

	out(w, value)

	close()
}
