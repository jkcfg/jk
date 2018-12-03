import flatbuffers from 'flatbuffers';
import { __std as fs } from '__std_FileSystem_generated';
import { __std as error } from '__std_Error_generated';
import { __std as std } from '__std_generated';

class FileInfo {
  constructor(p, d) {
    this.path = p;
    this.isdir = d;
  }
}

class Directory {
  constructor(p, files) {
    this.path = p;
    this.files = files;
  }
}

function info(path) {
  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
  fs.FileInfoArgs.startFileInfoArgs(builder);
  fs.FileInfoArgs.addPath(builder, pathOffset);
  const argsOffset = fs.FileInfoArgs.endFileInfoArgs(builder);

  std.Message.startMessage(builder);
  std.Message.addArgsType(builder, std.Args.FileInfoArgs);
  std.Message.addArgs(builder, argsOffset);
  const messageOffset = std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const bytes = V8Worker2.send(builder.asArrayBuffer());
  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = fs.FileSystemResponse.getRootAsFileSystemResponse(buf);
  switch (resp.retvalType()) {
  case fs.FileSystemRetval.FileInfo: {
    const f = new fs.FileInfo();
    resp.retval(f);
    return new FileInfo(f.path(), f.isdir());
  }
  case fs.FileSystemRetval.Error: {
    const err = new error.Error();
    resp.retval(err);
    throw new Error(err.message());
  }
  default:
    throw new Error('Unexpected response to fileinfo');
  }
}

function dir(path) {
  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
  fs.ListArgs.startListArgs(builder);
  fs.ListArgs.addPath(builder, pathOffset);
  const argsOffset = fs.ListArgs.endListArgs(builder);

  std.Message.startMessage(builder);
  std.Message.addArgsType(builder, std.Args.ListArgs);
  std.Message.addArgs(builder, argsOffset);
  const messageOffset = std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const bytes = V8Worker2.send(builder.asArrayBuffer());
  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = fs.FileSystemResponse.getRootAsFileSystemResponse(buf);
  switch (resp.retvalType()) {
  case fs.FileSystemRetval.Directory: {
    const d = new fs.Directory();
    resp.retval(d);
    const files = new Array(d.filesLength());
    for (let i = 0; i < files.length; i += 1) {
      const f = d.files(i);
      files[i] = new FileInfo(f.path(), f.isdir());
    }
    return new Directory(d.path(), files);
  }
  case fs.FileSystemRetval.Error: {
    const err = new error.Error();
    resp.retval(err);
    throw new Error(err.message());
  }
  default:
    throw new Error('Unexpected response to fileinfo');
  }
}

export {
  info,
  dir,
};
