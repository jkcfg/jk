package std

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jkcfg/jk/pkg/__std"
)

func param(params Params, kind __std.ParamType, path string, defaultValue string) []byte {
	var v interface{}
	var err error

	switch kind {
	case __std.ParamTypeBoolean:
		v, err = params.GetBool(path)
	case __std.ParamTypeNumber:
		v, err = params.GetNumber(path)
	case __std.ParamTypeString:
		v, err = params.GetString(path)
	case __std.ParamTypeObject:
		// For object parameters, we merge the default value with what the user has
		// given us, which could be only a part of the default value.
		v, err = params.GetObject(path)
		if err != nil {
			break
		}
		r := strings.NewReader(defaultValue)
		// The JS side sends us JSON.
		base, _ := NewParamsFromJSON(r)
		base.Merge(v.(Params))
		v = base
	default:
		panic("param: unexpected kind")
	}

	if err != nil && strings.Contains(err.Error(), "cannot convert") {
		// TODO(dlespiau): return an error to throw a JS exception.
		fmt.Fprintf(os.Stderr, "invalid type for param '%s': %v\n", path, err)
		return []byte("null")
	} else if err != nil {
		// path not found.
		return []byte("null")
	}

	// Param returns values that can be marshalled to JSON
	data, _ := json.Marshal(v)
	return data
}
