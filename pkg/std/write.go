package std

import (
	"encoding/json"
	"fmt"
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

type writerFunc func(io.Writer, []byte, int) error

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
func writeJSONFull(w io.Writer, value []byte, indent int, str writeString) error {
	var v interface{}
	if err := json.Unmarshal(value, &v); err != nil {
		return fmt.Errorf("writeJSON: unmarshal: %s", err.Error())
	}
	// Special case strings: we don't want to print them as JSON values.
	if s, ok := v.(string); str == rawString && ok {
		w.Write([]byte(s))
	} else {
		i, err := json.MarshalIndent(v, "", strings.Repeat(" ", indent))
		if err != nil {
			return fmt.Errorf("writeJSON: marshal: %s", err.Error())
		}
		w.Write(i)
	}
	_, err := w.Write([]byte{'\n'})
	return err
}

func writeJSON(str writeString) writerFunc {
	return func(w io.Writer, value []byte, indent int) error {
		return writeJSONFull(w, value, indent, str)
	}
}

func writeYAML(w io.Writer, value []byte, indent int) error {
	y, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return fmt.Errorf("writeYAML: %s", err.Error())
	}
	_, err = w.Write(y)
	return err
}

func writeYAMLStream(w io.Writer, v []byte, indent int) error {
	var values []interface{}
	if err := json.Unmarshal(v, &values); err != nil {
		return fmt.Errorf("writeYAMLStream: %s", err.Error())
	}
	encoder := yamlclassic.NewEncoder(w)
	for _, item := range values {
		if err := encoder.Encode(item); err != nil {
			return fmt.Errorf("writeYAMLStream: %s", err.Error())
		}
	}
	return nil
}

func writeJSONStream(w io.Writer, v []byte, indent int) error {
	var values []interface{}
	if err := json.Unmarshal(v, &values); err != nil {
		return fmt.Errorf("writeJSONStream: %s", err.Error())
	}
	encoder := json.NewEncoder(w)
	for _, item := range values {
		if err := encoder.Encode(item); err != nil {
			return fmt.Errorf("writeJSONStream: %s", err.Error())
		}
	}
	return nil
}

func writeHCL(w io.Writer, v []byte, indent int) error {
	ast, err := parser.Parse(v)
	if err != nil {
		return fmt.Errorf("writeHCL: unable to parse JSON: %s", err.Error())
	}

	config := printer.Config{
		SpacesWidth: indent,
	}
	err = config.Fprint(w, ast)
	if err != nil {
		return fmt.Errorf("writeHCL: unable to format HCL: %s", err.Error())
	}
	_, err = w.Write([]byte{'\n'})
	return err
}

func writeRaw(w io.Writer, value []byte, _ int) error {
	_, err := w.Write(value)
	return err
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

func write(value []byte, path string, format __std.Format, indent int, overwrite __std.Overwrite) error {
	switch overwrite {
	case __std.OverwriteWrite:
		break
	case __std.OverwriteSkip:
		if exists(path) {
			return nil
		}
	case __std.OverwriteErr:
		if exists(path) {
			return fmt.Errorf("file %s already exists", path)
		}
	}

	w, close := writer(path)

	var out writerFunc
	switch format {
	case __std.FormatFromExtension:
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
		return fmt.Errorf("write: unknown output format (%d)", int(format))
	}

	defer close()
	if err := out(w, value, indent); err != nil {
		return err
	}
	return nil
}

func (s Sandbox) Write(value []byte, path string, format __std.Format, indent int, overwrite __std.Overwrite) error {
	p, err := s.getWritePath(path)
	if err != nil {
		return err
	}
	return write(value, p, format, indent, overwrite)
}
