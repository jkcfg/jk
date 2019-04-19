package std

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jkcfg/jk/pkg/__std"
	"github.com/jkcfg/jk/pkg/__std/lib"
	"github.com/jkcfg/jk/pkg/deferred"

	flatbuffers "github.com/google/flatbuffers/go"
)

// Module returns the std source corresponding to the 'std' import
func Module(path string) []byte {
	data, _ := lib.ReadAll(path)
	return data
}

type sender interface {
	SendBytes([]byte) error
}

// ExecuteOptions global input parameters to the standards library.
type ExecuteOptions struct {
	// Verbose indicates if some operations, such as write, should print out what
	// they are doing.
	Verbose bool
	// Parameters is a structured set of input parameters.
	Parameters Params
	// OutputDirectory is a directory used by any file producing functions as the
	// base directory to output files to.
	OutputDirectory string
	// Root is topmost directory under which file reads are allowed
	Root ReadBase
}

func toBool(b byte) bool {
	if b == 0 {
		return false
	}
	return true
}

// stdError builds an Error flatbuffer we can return to the javascript side.
func stdError(b *flatbuffers.Builder, err error) flatbuffers.UOffsetT {
	off := b.CreateString(err.Error())
	__std.ErrorStart(b)
	__std.ErrorAddMessage(b, off)
	return __std.ErrorEnd(b)
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

		if path != "" && options.Verbose {
			fmt.Printf("wrote %s\n", path)
		}
		write(args.Value(), path, args.Format(), int(args.Indent()), toBool(args.Overwrite()))
		return nil
	case __std.ArgsReadArgs:
		args := __std.ReadArgs{}
		args.Init(union.Bytes, union.Pos)
		path := string(args.Path())
		if path != "" && options.Verbose {
			fmt.Printf("read %s\n", path)
		}
		module := string(args.Module())
		ser := deferred.Register(func() ([]byte, error) { return options.Root.Read(path, args.Format(), args.Encoding(), module) }, sendFunc(res.SendBytes))
		return deferredResponse(ser)

	case __std.ArgsFileInfoArgs:
		args := __std.FileInfoArgs{}
		args.Init(union.Bytes, union.Pos)
		return options.Root.FileInfo(string(args.Path()), string(args.Module()))
	case __std.ArgsListArgs:
		args := __std.ListArgs{}
		args.Init(union.Bytes, union.Pos)
		return options.Root.DirectoryListing(string(args.Path()), string(args.Module()))

	case __std.ArgsParamArgs:
		args := __std.ParamArgs{}
		args.Init(union.Bytes, union.Pos)

		// return buffer.
		b := flatbuffers.NewBuilder(512)
		var (
			off  flatbuffers.UOffsetT
			kind byte
		)

		json, err := param(options.Parameters, __std.ParamType(args.Type()), string(args.Path()), string(args.DefaultValue()))
		if err != nil {
			kind = __std.ParamRetvalError
			off = stdError(b, err)
		} else {
			kind = __std.ParamRetvalParamValue
			jsonOffset := b.CreateString(string(json))
			__std.ParamValueStart(b)
			__std.ParamValueAddJson(b, jsonOffset)
			off = __std.ParamValueEnd(b)
		}

		__std.ParamResponseStart(b)
		__std.ParamResponseAddRetvalType(b, kind)
		__std.ParamResponseAddRetval(b, off)
		responseOffset := __std.ParamResponseEnd(b)
		b.Finish(responseOffset)
		return b.FinishedBytes()

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
	off := stdError(b, err)
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
