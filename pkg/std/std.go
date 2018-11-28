package std

import (
	"log"
	"path/filepath"

	"github.com/justkidding-config/jk/pkg/__std"
	"github.com/justkidding-config/jk/pkg/__std/lib"
	"github.com/justkidding-config/jk/pkg/deferred"

	flatbuffers "github.com/google/flatbuffers/go"
)

// Module returns the std source corresponding to the 'std' import
func Module() []byte {
	data, _ := lib.ReadAll("std.js")
	return data
}

type sender interface {
	SendBytes([]byte) error
}

// ExecuteOptions global input parameters to the standards library.
type ExecuteOptions struct {
	// OutputDirectory is a directory used by any file producing functions as the
	// base directory to output files to.
	OutputDirectory string
}

// Execute parses a message from v8 and execute the corresponding function.
func Execute(msg []byte, res sender, options ExecuteOptions) []byte {
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
	case __std.ArgsReadArgs:
		args := __std.ReadArgs{}
		args.Init(union.Bytes, union.Pos)

		// TODO(michael): should do some validation and return an
		// error here when we handle more than one kind of thing, but
		// for now, treat everything as a file read from a local path
		// (which will only fail in the resolution, and can't be
		// cancelled).
		ser := deferred.Register(func() ([]byte, error) { return read(string(args.Url())) }, sendFunc(func(b []byte) error {
			return res.SendBytes(b)
		}))
		return deferredResponse(ser)
	default:
		log.Fatalf("unknown Message (%d)", message.ArgsType())
		return nil
	}
}

func deferredResponse(s deferred.Serial) []byte {
	b := flatbuffers.NewBuilder(20)
	__std.DeferredStart(b)
	__std.DeferredAddSerial(b, uint64(s))
	off := __std.DeferredEnd(b)
	__std.ResponseStart(b)
	__std.ResponseAddRetvalType(b, __std.RetvalDeferred)
	__std.ResponseAddRetval(b, off)
	off = __std.ResponseEnd(b)
	b.Finish(off)
	return b.FinishedBytes()
}

type sendFunc func([]byte) error

func (fn sendFunc) Error(s deferred.Serial, err error) {
	b := flatbuffers.NewBuilder(512)
	off := b.CreateString(err.Error())
	__std.ErrorStart(b)
	__std.ErrorAddMessage(b, off)
	off = __std.ErrorEnd(b)
	__std.ResolutionStart(b)
	__std.ResolutionAddSerial(b, uint64(s))
	__std.ResolutionAddValueType(b, __std.ResolutionValueError)
	__std.ResolutionAddValue(b, off)
	off = __std.ResolutionEnd(b)
	b.Finish(off)
	if err := fn(b.FinishedBytes()); err != nil {
		panic(err)
	}
}

func (fn sendFunc) Data(s deferred.Serial, data []byte) {
	b := flatbuffers.NewBuilder(1024)
	off := b.CreateByteString(data)
	__std.DataStart(b)
	__std.DataAddBytes(b, off)
	off = __std.DataEnd(b)
	__std.ResolutionStart(b)
	__std.ResolutionAddSerial(b, uint64(s))
	__std.ResolutionAddValueType(b, __std.ResolutionValueData)
	__std.ResolutionAddValue(b, off)
	off = __std.ResolutionEnd(b)
	b.Finish(off)
	if err := fn(b.FinishedBytes()); err != nil {
		panic(err)
	}
}

func (fn sendFunc) End(s deferred.Serial) {
	b := flatbuffers.NewBuilder(1024)
	__std.EndOfStreamStart(b)
	off := __std.EndOfStreamEnd(b)
	__std.ResolutionStart(b)
	__std.ResolutionAddSerial(b, uint64(s))
	__std.ResolutionAddValueType(b, __std.ResolutionValueEndOfStream)
	__std.ResolutionAddValue(b, off)
	off = __std.ResolutionEnd(b)
	b.Finish(off)
	if err := fn(b.FinishedBytes()); err != nil {
		panic(err)
	}
}
