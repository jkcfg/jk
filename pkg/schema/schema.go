// Package schema is an adapter from
// https://github.com/xeipuuv/gojsonschema to the jk runtime, so that
// validation via JSON Schema can be built into the standard library.
package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
)

// ValidateWithObject validates a value against a schema. Both
// arguments are supplied as strings, since they will be passed as
// strings via RPC anyway.
func ValidateWithObject(valueStr, schemaStr string) ([]string, error) {
	valueLoader := jsonschema.NewStringLoader(valueStr)
	schemaLoader := jsonschema.NewStringLoader(schemaStr)

	result, err := jsonschema.Validate(schemaLoader, valueLoader)
	if err != nil {
		return nil, err
	}

	var errors []string
	for _, result := range result.Errors() {
		errors = append(errors, result.Description())
	}
	return errors, nil
}
