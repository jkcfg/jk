package std

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func p(s string) Params {
	r := strings.NewReader(s)
	params, err := NewParamsFromJSON(r)
	if err != nil {
		log.Fatal(s, err)
	}
	return params
}

func TestGet(t *testing.T) {
	tests := []struct {
		o        Params
		path     string
		valid    bool
		expected interface{}
	}{
		{map[string]interface{}{}, "foo.bar", false, nil},
		{p(`{ "foo": 2 }`), "foo.bar", false, nil},
		{p(`{ "foo": { "bar": 2 } }`), "foo.bar", true, float64(2)},
		{p(`{ "foo": { "bar": "baz" } }`), "foo.bar", true, "baz"},
		{p(`{ "foo": { "bar": { "baz": 3 } } }`), "foo.bar", true, p(`{ "baz": 3 }`)},
		{p(`{ "xxx": "yyy", "foo": { "bar": { "baz": 3 } } }`), "foo.bar", true, p(`{ "baz": 3 }`)},
	}

	for _, test := range tests {
		v, err := test.o.Get(test.path)
		if test.valid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
		assert.Equal(t, test.expected, v)
	}
}

func TestTypedGet(t *testing.T) {
	params := p(`{ "xxx": "yyy", "foo": { "bar": { "baz": 3 }, "boolean": true } }`)

	vBool, err := params.GetBool("foo.boolean")
	assert.NoError(t, err)
	assert.Equal(t, true, vBool)

	vNumber, err := params.GetNumber("foo.bar.baz")
	assert.NoError(t, err)
	assert.Equal(t, float64(3), vNumber)

	vString, err := params.GetString("xxx")
	assert.NoError(t, err)
	assert.Equal(t, "yyy", vString)
	_, err = params.GetString("foo.bar.baz")
	assert.Error(t, err)

	vObject, err := params.GetObject("foo.bar")
	assert.NoError(t, err)
	assert.Equal(t, p(`{ "baz": 3 }`), vObject)

}

func TestSet(t *testing.T) {
	tests := []struct {
		o        Params
		path     string
		value    interface{}
		expected Params
	}{
		{p(`{}`), "foo", float64(2), p(`{ "foo": 2 }`)},
		{p(`{}`), "foo", "bar", p(`{ "foo": "bar" }`)},
		{p(`{}`), "foo", true, p(`{ "foo": true }`)},
		{p(`{}`), "foo", p(`{ "bar": "baz" } `), p(`{ "foo": { "bar": "baz" } }`)},
		{p(`{ "foo": { "xxx": 42 } }`), "foo.yyy", p(`{ "bar": "baz" } `), p(`{ "foo": { "xxx": 42, "yyy": { "bar": "baz" } } }`)},
	}

	for _, test := range tests {
		test.o.Set(test.path, test.value)
		assert.Equal(t, test.expected, test.o)
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		a, b     Params
		expected Params
	}{
		{NewParams(), NewParams(), NewParams()},
		{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}},
		{map[string]interface{}{}, map[string]interface{}{"foo": 1}, map[string]interface{}{"foo": 1}},
		{map[string]interface{}{}, map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"}},
		{map[string]interface{}{"foo": 1}, map[string]interface{}{}, map[string]interface{}{"foo": 1}},
		{map[string]interface{}{"foo": "bar"}, map[string]interface{}{}, map[string]interface{}{"foo": "bar"}},

		{p(`{ "foo": 1 } `), p(`{ "foo": { "bar": "baz" } }`), p(`{"foo": { "bar": "baz" } }`)},
		{p(`{ "foo": 1, "orig": "xxx" } `), p(`{ "foo": { "bar": "baz" } }`), p(`{"foo": { "bar": "baz" }, "orig": "xxx" }`)},
		{p(`{ "foo": { "rab": "zab" }, "orig": "xxx" } `), p(`{ "foo": { "bar": "baz" } }`), p(`{"foo": { "bar": "baz", "rab": "zab" }, "orig": "xxx" }`)},
		{p(`{ "foo": { "bar": "baz" } }`), p(`{ "foo": { "rab": "zab" }, "orig": "xxx" } `), p(`{"foo": { "bar": "baz", "rab": "zab" }, "orig": "xxx" }`)},
	}

	for _, test := range tests {
		test.a.Merge(test.b)
		assert.Equal(t, test.expected, test.a)
	}

}
