package std

import (
	"log"
	"path/filepath"

	"github.com/justkidding-config/jk/pkg/__std"
	"github.com/justkidding-config/jk/pkg/__std/lib"

	flatbuffers "github.com/google/flatbuffers/go"
)

// Module returns the std source corresponding to the 'std' import
func Module() []byte {
	data, _ := lib.ReadAll("std.js")
	return data
}

// ExecuteOptions global input parameters to the standards library.
type ExecuteOptions struct {
	// OutputDirectory is a directory used by any file producing functions as the
	// base directory to output files to.
	OutputDirectory string
}

// Execute parses a message from v8 and execute the corresponding function.
func Execute(msg []byte, options ExecuteOptions) []byte {
	message := __std.GetRootAsMessage(msg, 0)

	union := flatbuffers.Table{}
	if !message.Args(&union) {
		log.Fatal("could not parse Message")
	}

	switch message.ArgsType() {
	case __std.ArgsWriteArgs:
		args := __std.WriteArgs{}
		args.Init(union.Bytes, union.Pos)

		// Weave options.OutputDirectory in there.
		path := string(args.Path())
		if path != "" && !filepath.IsAbs(path) {
			path = filepath.Join(options.OutputDirectory, path)
		}

		write(args.Value(), path, args.Type(), int(args.Indent()))
		return nil
	default:
		log.Fatalf("unknown Message (%d)", message.ArgsType())
		return nil
	}
}
