// Package schema is an adapter from
// https://github.com/xeipuuv/gojsonschema to the jk runtime, so that
// validation via JSON Schema can be built into the standard library.
package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
)

// Error is a struct that will serialise as JSON neatly, to be
// rehydrated as a validation error by the RPC wrapper in JavaScript.
type Error struct {
	Msg  string `json:"msg"`
	Path string `json:"path"`
}

func validate(valueLoader, schemaLoader jsonschema.JSONLoader) ([]Error, error) {
	result, err := jsonschema.Validate(schemaLoader, valueLoader)
	if err != nil {
		return nil, err
	}

	var errors []Error
	for _, result := range result.Errors() {
		errors = append(errors, Error{
			Msg:  result.Description(),
			Path: result.Field(),
		})
	}
	return errors, nil
}

// ValidateWithObject validates a value against a schema. Both
// arguments are supplied as strings, since they will be passed as
// strings via RPC anyway.
func ValidateWithObject(valueStr, schemaStr string) ([]Error, error) {
	valueLoader := jsonschema.NewStringLoader(valueStr)
	schemaLoader := jsonschema.NewStringLoader(schemaStr)
	return validate(valueLoader, schemaLoader)
}

// ValidateWithFile validates a value (as JSON stringified) against
// the schema at the path given.
func ValidateWithFile(valueStr, path string) ([]Error, error) {
	valueLoader := jsonschema.NewStringLoader(valueStr)
	schemaLoader := jsonschema.NewReferenceLoader("file://" + path)
	return validate(valueLoader, schemaLoader)
}
