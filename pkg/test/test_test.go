package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	tests := []struct {
		test     *Test
		expected string
	}{
		{New("test-foo.js"), "foo"},
		{New("foo.js"), "foo"},
		{New("foo.js", Options{Name: "bar"}), "bar"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.test.Name())
	}
}
