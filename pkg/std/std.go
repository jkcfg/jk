package std

import (
	"log"

	"github.com/dlespiau/jk/pkg/__std"
	"github.com/dlespiau/jk/pkg/__std/lib"

	flatbuffers "github.com/google/flatbuffers/go"
)

// Module returns the std source corresponding to the 'std' import
func Module() []byte {
	data, _ := lib.ReadAll("std.js")
	return data
}

// Execute parses a message from v8 and execute the corresponding function.
func Execute(msg []byte) []byte {
	message := __std.GetRootAsMessage(msg, 0)

	union := flatbuffers.Table{}
	if !message.Args(&union) {
		log.Fatal("could not parse Message")
	}

	switch message.ArgsType() {
	case __std.ArgsWriteArgs:
		args := __std.WriteArgs{}
		args.Init(union.Bytes, union.Pos)
		write(args.Value(), string(args.Path()), args.Type())
		return nil
	default:
		log.Fatalf("unknown Message (%d)", message.ArgsType())
		return nil
	}
}
