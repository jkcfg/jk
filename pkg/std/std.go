package std

import (
	"context"
	"log"
	"path/filepath"

	"github.com/jkcfg/jk/pkg/__std"
	"github.com/jkcfg/jk/pkg/__std/lib"
	"github.com/jkcfg/jk/pkg/deferred"

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
	case __std.ArgsCancelArgs:
		args := __std.CancelArgs{}
		args.Init(union.Bytes, union.Pos)
		serial := args.Serial()
		deferred.Cancel(deferred.Serial(serial))
		return nil
	case __std.ArgsReadArgs:
		args := __std.ReadArgs{}
		args.Init(union.Bytes, union.Pos)

		// TODO(michael): should do some validation and return an
		// error here when we handle more than one kind of thing, but
		// for now, treat everything as a file read from a local path
		// (which will only fail in the resolution, and can't be
		// cancelled).
		ser := deferred.Register(func(_ context.Context) ([]byte, error) { return read(string(args.Url())) }, sendFunc(res.SendBytes))
		return deferredResponse(ser)
	case __std.ArgsFileInfoArgs:
		args := __std.FileInfoArgs{}
		args.Init(union.Bytes, union.Pos)
		return fileInfo(string(args.Path()))
	case __std.ArgsListArgs:
		args := __std.ListArgs{}
		args.Init(union.Bytes, union.Pos)
		return directoryListing(string(args.Path()))
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
	__std.DeferredResponseStart(b)
	__std.DeferredResponseAddRetvalType(b, __std.DeferredRetvalDeferred)
	__std.DeferredResponseAddRetval(b, off)
	off = __std.DeferredResponseEnd(b)
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
	__std.FulfilmentStart(b)
	__std.FulfilmentAddSerial(b, uint64(s))
	__std.FulfilmentAddValueType(b, __std.FulfilmentValueError)
	__std.FulfilmentAddValue(b, off)
	off = __std.FulfilmentEnd(b)
	b.Finish(off)
	if err := fn(b.FinishedBytes()); err != nil {
		panic(err)
	}
}

func (fn sendFunc) Data(s deferred.Serial, data []byte) {
	b := flatbuffers.NewBuilder(1024)
	off := b.CreateByteVector(data)
	__std.DataStart(b)
	__std.DataAddBytes(b, off)
	off = __std.DataEnd(b)
	__std.FulfilmentStart(b)
	__std.FulfilmentAddSerial(b, uint64(s))
	__std.FulfilmentAddValueType(b, __std.FulfilmentValueData)
	__std.FulfilmentAddValue(b, off)
	off = __std.FulfilmentEnd(b)
	b.Finish(off)
	if err := fn(b.FinishedBytes()); err != nil {
		panic(err)
	}
}

func (fn sendFunc) End(s deferred.Serial) {
	b := flatbuffers.NewBuilder(1024)
	__std.EndOfStreamStart(b)
	off := __std.EndOfStreamEnd(b)
	__std.FulfilmentStart(b)
	__std.FulfilmentAddSerial(b, uint64(s))
	__std.FulfilmentAddValueType(b, __std.FulfilmentValueEndOfStream)
	__std.FulfilmentAddValue(b, off)
	off = __std.FulfilmentEnd(b)
	b.Finish(off)
	if err := fn(b.FinishedBytes()); err != nil {
		panic(err)
	}
}
