package std

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Params is a paramater store akin to a JSON object.
type Params map[string]interface{}

// NewParams creates an empty set of parameters.
func NewParams() Params {
	return make(map[string]interface{})
}

// NewParamsFromJSON creates Params from JSON.
func NewParamsFromJSON(r io.Reader) (Params, error) {
	p := NewParams()
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&p)
	return p, err
}

// Get retrieves a parameter.
func (p Params) Get(path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	m := p
	for i, part := range parts {
		v, found := m[part]
		if !found {
			return nil, fmt.Errorf("invalid path (key not found): %s", strings.Join(parts[:i+1], "."))
		}
		// Found the value!
		if i == len(parts)-1 {
			// When the value is a sub-object (as opposed to a primitive value), we box
			// it into a Params to keep the nice property that a sub-tree of Params is
			// still a Params.
			if retMap, ok := v.(map[string]interface{}); ok {
				return Params(retMap), nil
			}
			return v, nil
		}
		// We can only continue if we're traversing a map.
		if newMap, ok := v.(map[string]interface{}); ok {
			m = Params(newMap)
			continue
		}
		return nil, fmt.Errorf("invalid path (key isn't a map): %s", strings.Join(parts[:i+1], "."))
	}

	// We shouldn't reach this.
	return nil, fmt.Errorf("invalid path (eek!): %s", path)
}

// GetBool retrieves a boolean parameter.
func (p Params) GetBool(path string) (bool, error) {
	v, err := p.Get(path)
	if err != nil {
		return false, err
	}
	if b, ok := v.(bool); ok {
		return b, nil
	}
	// string -> bool coercion
	if s, ok := v.(string); ok {
		if b, err := strconv.ParseBool(s); err == nil {
			return b, nil
		}
	}
	return false, fmt.Errorf("cannot convert %v to bool", v)
}

// GetNumber retrieves a number parameter.
func (p Params) GetNumber(path string) (float64, error) {
	v, err := p.Get(path)
	if err != nil {
		return 0, err
	}
	if f, ok := v.(float64); ok {
		return f, nil
	}
	// string -> number coercion.
	if s, ok := v.(string); ok {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f, nil
		}
	}
	return 0, fmt.Errorf("cannot convert %v to float64", v)
}

// GetString retrieves a string parameter.
func (p Params) GetString(path string) (string, error) {
	v, err := p.Get(path)
	if err != nil {
		return "", err
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("cannot convert %v to string", v)
}

// GetObject retrieves a object parameter.
func (p Params) GetObject(path string) (Params, error) {
	v, err := p.Get(path)
	if err != nil {
		return NewParams(), err
	}
	if o, ok := v.(Params); ok {
		return o, nil
	}
	return NewParams(), fmt.Errorf("cannot convert %v to Params", v)
}

func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

// Set sets a parameter.
func (p Params) Set(path string, v interface{}) {
	parts := strings.Split(path, ".")
	m := p
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		value, ok := m[part]
		// If the part key doesn't exist yet or is primitive type, create it.
		if !ok || !isMap(value) {
			p := make(map[string]interface{})
			m[part] = p
			m = p
			continue
		}
		// Continue through the chain of maps
		m = m[part].(map[string]interface{})
	}

	// Internal types are only primitive types and maps, not Params. Convert Params
	// to map[string]interface{}.
	if p, ok := v.(Params); ok {
		v = map[string]interface{}(p)
	}

	key := parts[len(parts)-1]
	m[key] = v
}

// SetBool sets a boolean parameter.
func (p Params) SetBool(path string, b bool) {
	p.Set(path, b)
}

// SetNumber sets a number parameter.
func (p Params) SetNumber(path string, f float64) {
	p.Set(path, f)
}

// SetString sets a string parameter.
func (p Params) SetString(path string, s string) {
	p.Set(path, s)
}

// SetObject sets an object parameter.
func (p Params) SetObject(path string, o Params) {
	p.Set(path, o)
}

// Merge merges two parameter stores.
func (p Params) Merge(a Params) {
	for k, v := range a {
		// k isn't in the original map, set it.
		if _, ok := p[k]; !ok {
			p[k] = v
			continue
		}
		dst := p[k]

		// Both original and incoming values are maps.
		dstMap, dstOk := dst.(map[string]interface{})
		srcMap, srcOk := v.(map[string]interface{})
		if dstOk && srcOk {
			Params(dstMap).Merge(srcMap)
			continue
		}

		// Otherwise, source overrides destination
		p[k] = v
	}
}

func (p Params) String() string {
	s, _ := json.MarshalIndent(p, "", " ")
	return string(s)
}
