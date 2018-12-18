import flatbuffers from 'flatbuffers';
import { __std } from '__std_generated';

class FileInfo {
  constructor(n, p, d) {
    this.name = n;
    this.path = p;
    this.isdir = d;
  }
}

class Directory {
  constructor(n, p, files) {
    this.name = n;
    this.path = p;
    this.files = files;
  }
}

function info(path) {
  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
  __std.FileInfoArgs.startFileInfoArgs(builder);
  __std.FileInfoArgs.addPath(builder, pathOffset);
  const argsOffset = __std.FileInfoArgs.endFileInfoArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.FileInfoArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const bytes = V8Worker2.send(builder.asArrayBuffer());
  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = __std.FileSystemResponse.getRootAsFileSystemResponse(buf);
  switch (resp.retvalType()) {
  case __std.FileSystemRetval.FileInfo: {
    const f = new __std.FileInfo();
    resp.retval(f);
    return new FileInfo(f.name(), f.path(), f.isdir());
  }
  case __std.FileSystemRetval.Error: {
    const err = new __std.Error();
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
  __std.ListArgs.startListArgs(builder);
  __std.ListArgs.addPath(builder, pathOffset);
  const argsOffset = __std.ListArgs.endListArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.ListArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const bytes = V8Worker2.send(builder.asArrayBuffer());
  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = __std.FileSystemResponse.getRootAsFileSystemResponse(buf);
  switch (resp.retvalType()) {
  case __std.FileSystemRetval.Directory: {
    const d = new __std.Directory();
    resp.retval(d);
    const files = new Array(d.filesLength());
    for (let i = 0; i < files.length; i += 1) {
      const f = d.files(i);
      files[i] = new FileInfo(f.name(), f.path(), f.isdir());
    }
    return new Directory(d.name(), d.path(), files);
  }
  case __std.FileSystemRetval.Error: {
    const err = new __std.Error();
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
