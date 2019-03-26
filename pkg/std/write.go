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
	yamlclassic "gopkg.in/yaml.v2"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/json/parser"
)

type closer func()

func nilCloser() {}

type writerFunc func(io.Writer, []byte, int)

type writeString bool

const (
	rawString  writeString = true
	jsonString writeString = false
)

// We have a special case: when we're asked to print a string, we just write it
// instead of writing the json value of a string.
//   std.log('foo') -> foo (not "foo")
// However, when explicitly asked to write JSON, we still need to honour that:
//   std.write('foo', 'file.json') -> "foo"
func writeJSONFull(w io.Writer, value []byte, indent int, str writeString) {
	var v interface{}
	if err := json.Unmarshal(value, &v); err != nil {
		log.Fatalf("writeJSON: unmarshal: %s", err)
	}
	// Special case strings: we don't want to print them as JSON values.
	if s, ok := v.(string); str == rawString && ok {
		w.Write([]byte(s))
	} else {
		i, err := json.MarshalIndent(v, "", strings.Repeat(" ", indent))
		if err != nil {
			log.Fatalf("writeJSON: marshal: %s", err)
		}
		w.Write(i)
	}
	w.Write([]byte{'\n'})
}

func writeJSON(str writeString) writerFunc {
	return func(w io.Writer, value []byte, indent int) {
		writeJSONFull(w, value, indent, str)
	}
}

func writeYAML(w io.Writer, value []byte, indent int) {
	y, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		log.Fatalf("writeYAML: %s", err)
	}
	w.Write(y)
}

func writeYAMLStream(w io.Writer, v []byte, indent int) {
	var values []interface{}
	if err := json.Unmarshal(v, &values); err != nil {
		log.Fatalf("writeYAMLStream: %s", err)
	}
	encoder := yamlclassic.NewEncoder(w)
	for _, item := range values {
		if err := encoder.Encode(item); err != nil {
			log.Fatalf("writeYAMLStream: %s", err)
		}
	}
}

func writeJSONStream(w io.Writer, v []byte, indent int) {
	var values []interface{}
	if err := json.Unmarshal(v, &values); err != nil {
		log.Fatalf("writeJSONStream: %s", err)
	}
	encoder := json.NewEncoder(w)
	for _, item := range values {
		if err := encoder.Encode(item); err != nil {
			log.Fatalf("writeJSONStream: %s", err)
		}
	}
}

func writeHCL(w io.Writer, v []byte, indent int) {
	ast, err := parser.Parse(v)
	if err != nil {
		log.Fatalf("writeHCL: unable to parse JSON: %s", err)
	}

	config := printer.Config{
		SpacesWidth: indent,
	}
	err = config.Fprint(w, ast)
	if err != nil {
		log.Fatalf("writeHCL: unable to format HCL: %s", err)
	}
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
		return writeJSON(jsonString)
	case ".hcl", ".tf":
		return writeHCL
	default:
		return writeJSON(rawString)
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func write(value []byte, path string, format __std.Format, indent int, overwrite bool) {
	if !overwrite && exists(path) {
		return
	}

	w, close := writer(path)

	var out writerFunc
	switch format {
	case __std.FormatAuto:
		out = writerFuncFromPath(path)
	case __std.FormatJSON:
		out = writeJSON(jsonString)
	case __std.FormatJSONStream:
		out = writeJSONStream
	case __std.FormatYAML:
		out = writeYAML
	case __std.FormatYAMLStream:
		out = writeYAMLStream
	case __std.FormatHCL:
		out = writeHCL
	case __std.FormatRaw:
		out = writeRaw
	default:
		log.Fatalf("write: unknown output format (%d)", int(format))
	}

	out(w, value, indent)

	close()
}
