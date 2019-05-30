package std

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jkcfg/jk/pkg/__std"

	"github.com/ghodss/yaml"
	flatbuffers "github.com/google/flatbuffers/go"
)

// Parse accepts a stringified object, and the format in which it was
// stringified, and returns a JSON stringified object.
func Parse(input []byte, format __std.Format) ([]byte, error) {
	switch format {
	case __std.FormatJSON:
		// This is strict: exactly one JSON-encoded object is allowed
		// (that's how json.Unmarshal works).
		var throwaway interface{}
		if err := json.Unmarshal(input, &throwaway); err != nil {
			return nil, err
		}
		return input, nil
	case __std.FormatYAML:
		return yaml.YAMLToJSON(input)
	case __std.FormatJSONStream:
		return readJSONStream(bytes.NewReader(input))
	case __std.FormatYAMLStream:
		return readYAMLStream(bytes.NewReader(input))
	}
	return nil, fmt.Errorf(`Unsupported format for Parse: %s`, __std.EnumNamesFormat[format])
}

// Unparse accepts a JSON stringified object, and a format for
// reserialising it, and returns the reserialised object.
func Unparse(jsonString []byte, format __std.Format) ([]byte, error) {
	var value interface{}
	if err := json.Unmarshal(jsonString, &value); err != nil {
		return nil, err
	}
	switch format {
	case __std.FormatJSON:
		return jsonString, nil
	}
	return nil, fmt.Errorf(`Unsupported format for Unparse: %s`, __std.EnumNamesFormat[format])
}

// parseUnparseReturn constructs a flatbuffer-encoded return value
// given bytes and an error (either of which may be nil).
func parseUnparseResponse(out []byte, err error) []byte {
	b := flatbuffers.NewBuilder(1024)

	if err != nil {
		off := stdError(b, err)
		__std.ParseUnparseResponseStart(b)
		__std.ParseUnparseResponseAddRetvalType(b, __std.ParseUnparseRetvalError)
		__std.ParseUnparseResponseAddRetval(b, off)
		off = __std.ParseUnparseResponseEnd(b)
		b.Finish(off)
		return b.FinishedBytes()
	}

	off := b.CreateByteString(out)
	__std.ParseUnparseDataStart(b)
	__std.ParseUnparseDataAddData(b, off)
	off = __std.ParseUnparseDataEnd(b)
	__std.ParseUnparseResponseStart(b)
	__std.ParseUnparseResponseAddRetvalType(b, __std.ParseUnparseRetvalParseUnparseData)
	__std.ParseUnparseResponseAddRetval(b, off)
	off = __std.ParseUnparseResponseEnd(b)
	b.Finish(off)
	return b.FinishedBytes()
}
