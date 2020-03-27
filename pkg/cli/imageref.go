package cli

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

// ImageRefSliceValue is a slice of image references (name.Reference) for
// use with pflag. Adapted from
// https://github.com/spf13/pflag/blob/master/string_slice.go
type ImageRefSliceValue struct {
	refs    *[]name.Reference
	changed bool
}

// NewImageRefSliceValue creates a `pflag.Value` wrapper for a slice
// of name.Reference (image reference) values.
func NewImageRefSliceValue(p *[]name.Reference) *ImageRefSliceValue {
	v := new(ImageRefSliceValue)
	v.refs = p
	return v
}

func writeRefsAsCsv(vals []name.Reference) (string, error) {
	strs := make([]string, len(vals), len(vals))
	for i := range vals {
		strs[i] = vals[i].String()
	}
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(strs)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

// String returns a string representation of the value
func (value *ImageRefSliceValue) String() string {
	str, _ := writeRefsAsCsv(*value.refs)
	return "[" + str + "]"
}

// Set sets the underlying value given a string passed to the flag. As
// with other slice values, if it has already been set, being called
// again will append to the slice.
func (value *ImageRefSliceValue) Set(v string) error {
	strs, err := readAsCSV(v)
	if err != nil && err != io.EOF {
		return err
	}

	out := make([]name.Reference, len(strs), len(strs))
	for i := range strs {
		r, err := name.ParseReference(strs[i])
		if err != nil {
			return err
		}
		// we want either a tag or a digest
		if r.Identifier() == "latest" {
			return fmt.Errorf("image ref has no tag or digest, or uses 'latest'")
		}

		out[i] = r
	}

	if !value.changed {
		*value.refs = out
	} else {
		*value.refs = append(*value.refs, out...)
	}

	value.changed = true
	return nil
}

// Type returns a name for the type of flag
func (*ImageRefSliceValue) Type() string {
	return "imageRefSlice"
}

// getStringSlice is a handy way to get a slice of strings from the
// value, for testing.
func (value *ImageRefSliceValue) getStringSlice() []string {
	strs := make([]string, len(*value.refs), len(*value.refs))
	for i, ref := range *value.refs {
		strs[i] = ref.String()
	}
	return strs
}
