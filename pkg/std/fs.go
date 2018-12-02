package std

import (
	"os"

	"github.com/justkidding-config/jk/pkg/__std"

	flatbuffers "github.com/google/flatbuffers/go"
)

func fileInfo(path string) []byte {
	info, err := os.Stat(path)
	switch {
	case err != nil:
		return fsError(err.Error())
	case !(info.IsDir() || info.Mode().IsRegular()):
		return fsError("not a regular file")
	}
	return fileInfoResponse(path, info.IsDir())
}

func fsError(msg string) []byte {
	b := flatbuffers.NewBuilder(1024)
	off := b.CreateString(msg)
	__std.ErrorStart(b)
	__std.ErrorAddMessage(b, off)
	off = __std.ErrorEnd(b)
	__std.FileSystemResponseStart(b)
	__std.FileSystemResponseAddRetvalType(b, __std.FileSystemRetvalError)
	__std.FileSystemResponseAddRetval(b, off)
	off = __std.FileSystemResponseEnd(b)
	b.Finish(off)
	return b.FinishedBytes()
}

func fileInfoResponse(path string, isdir bool) []byte {
	b := flatbuffers.NewBuilder(1024)
	off := buildFileInfo(b, path, isdir)
	__std.FileSystemResponseStart(b)
	__std.FileSystemResponseAddRetvalType(b, __std.FileSystemRetvalFileInfo)
	__std.FileSystemResponseAddRetval(b, off)
	off = __std.FileSystemResponseEnd(b)
	b.Finish(off)
	return b.FinishedBytes()
}

func buildFileInfo(b *flatbuffers.Builder, path string, isdir bool) flatbuffers.UOffsetT {
	off := b.CreateString(path)
	__std.FileInfoStart(b)
	__std.FileInfoAddPath(b, off)
	if isdir {
		__std.FileInfoAddIsdir(b, 1)
	} else {
		__std.FileInfoAddIsdir(b, 0)
	}
	return __std.FileInfoEnd(b)
}

func directoryListing(path string) []byte {
	dir, err := os.Open(path)
	if err != nil {
		return fsError(err.Error())
	}
	infos, err := dir.Readdir(0)
	if err != nil {
		return fsError(err.Error())
	}

	b := flatbuffers.NewBuilder(1024)
	offsets := make([]flatbuffers.UOffsetT, len(infos), len(infos))
	for i := range infos {
		offsets[i] = buildFileInfo(b, infos[i].Name(), infos[i].IsDir())
	}

	__std.DirectoryStartFilesVector(b, len(offsets))
	for i := len(offsets) - 1; i >= 0; i-- {
		b.PrependUOffsetT(offsets[i])
	}
	infoVec := b.EndVector(len(offsets))

	pathOff := b.CreateString(path)
	__std.DirectoryStart(b)
	__std.DirectoryAddPath(b, pathOff)
	__std.DirectoryAddFiles(b, infoVec)
	dirOff := __std.DirectoryEnd(b)

	__std.FileSystemResponseStart(b)
	__std.FileSystemResponseAddRetvalType(b, __std.FileSystemRetvalDirectory)
	__std.FileSystemResponseAddRetval(b, dirOff)
	off := __std.FileSystemResponseEnd(b)
	b.Finish(off)
	return b.FinishedBytes()
}
