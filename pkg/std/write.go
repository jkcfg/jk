package std

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jkcfg/jk/pkg/__std"

	"github.com/ghodss/yaml"
)

type closer func()

func nilCloser() {}

type writerFunc func(io.Writer, []byte, int)

func writeJSON(w io.Writer, value []byte, indent int) {
	var v interface{}
	if err := json.Unmarshal(value, &v); err != nil {
		log.Fatalf("writeJSON: unmarshal: %s", err)
	}
	i, err := json.MarshalIndent(v, "", strings.Repeat(" ", indent))
	if err != nil {
		log.Fatalf("writeJSON: marshal: %s", err)
	}
	w.Write(i)
	w.Write([]byte{'\n'})
}

func writeYAML(w io.Writer, value []byte, indent int) {
	y, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		log.Fatalf("writeYAML: %s", err)
	}
	w.Write(y)
}

func writeRaw(w io.Writer, value []byte, _ int) {
	w.Write(value)
}

func writer(path string) (io.Writer, closer) {
	if path == "" {
		return os.Stdout, nilCloser
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0770); err != nil {
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

func write(value []byte, path string, format __std.Format, indent int) {
	w, close := writer(path)

	var out writerFunc
	switch format {
	case __std.FormatAuto:
		out = writerFuncFromPath(path)
	case __std.FormatJSON:
		out = writeJSON
	case __std.FormatYAML:
		out = writeYAML
	case __std.FormatRaw:
		out = writeRaw
	default:
		log.Fatalf("write: unknown output format (%d)", int(format))
	}

	out(w, value, indent)

	close()
}
