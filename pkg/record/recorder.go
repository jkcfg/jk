package record

import "encoding/json"

// OperationKind is the name of a recorded operation
type OperationKind string

// Params are Operation parameters we can record along with an entry.
type Params map[string]interface{}

const (
	// ImportFile is a js import from the filesystem (exclude stdlib imports).
	ImportFile OperationKind = "import-file"
	// ReadFile is a std.read from the filesystem (exclude reading from stdin).
	ReadFile OperationKind = "read-file"
)

// Operation is an entry in the Recording.
type Operation struct {
	kind   OperationKind
	params Params
}

// MarshalJSON implements json.Marshaler.
func (o Operation) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["kind"] = o.kind
	for k, v := range o.params {
		m[k] = v
	}
	return json.Marshal(m)
}

// Recorder records Operations. It's designed to be generic, any part of jk can
// append an operation to the log.
type Recorder struct {
	ops []Operation
}

// Record appends a new operation to the log.
func (r *Recorder) Record(kind OperationKind, params Params) {
	r.ops = append(r.ops, Operation{kind: kind, params: params})
}

// Log retrieves the recording list of operations
func (r *Recorder) Log() []Operation {
	return r.ops
}

// MarshalJSON implements json.Marshaler.
func (r Recorder) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ops)
}
