import { flatbuffers } from './internal/flatbuffers';
import { __std } from './internal/__std_generated';

class FileInfo {
  name: string;
  path: string;
  isdir: boolean;

  constructor(n: string, p: string, d: boolean) {
    this.name = n;
    this.path = p;
    this.isdir = d;
  }
}

class Directory {
  name: string;
  path: string;
  files: FileInfo[];

  constructor(n: string, p: string, files: FileInfo[]) {
    this.name = n;
    this.path = p;
    this.files = files;
  }
}

interface InfoOptions {
  module?: string;
}

function info(path: string, { module }: InfoOptions = {}): FileInfo {
  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
  let moduleOffset = 0;
  if (module !== undefined) {
    moduleOffset = builder.createString(module);
  }

  __std.FileInfoArgs.startFileInfoArgs(builder);
  __std.FileInfoArgs.addPath(builder, pathOffset);
  if (module !== undefined) {
    __std.FileInfoArgs.addModule(builder, moduleOffset);
  }
  const argsOffset = __std.FileInfoArgs.endFileInfoArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.FileInfoArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const bytes = <ArrayBuffer>V8Worker2.send(builder.asArrayBuffer());
  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = __std.FileSystemResponse.getRootAsFileSystemResponse(buf);
  switch (resp.retvalType()) {
  case __std.FileSystemRetval.FileInfo: {
    const f = new __std.FileInfo();
    resp.retval(f);
    return new FileInfo(<string>f.name(), <string>f.path(), f.isdir());
  }
  case __std.FileSystemRetval.Error: {
    const err = new __std.Error();
    resp.retval(err);
    throw new Error(<string>err.message());
  }
  default:
    throw new Error('Unexpected response to fileinfo');
  }
}

interface DirOptions {
  module?: string;
}

function dir(path: string, { module }: DirOptions = {}): Directory {
  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
  let moduleOffset = 0;
  if (module !== undefined) {
    moduleOffset = builder.createString(module);
  }

  __std.ListArgs.startListArgs(builder);
  __std.ListArgs.addPath(builder, pathOffset);
  if (module !== undefined) {
    __std.ListArgs.addModule(builder, moduleOffset);
  }
  const argsOffset = __std.ListArgs.endListArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.ListArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const bytes = <ArrayBuffer>V8Worker2.send(builder.asArrayBuffer());
  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = __std.FileSystemResponse.getRootAsFileSystemResponse(buf);
  switch (resp.retvalType()) {
  case __std.FileSystemRetval.Directory: {
    const d = new __std.Directory();
    resp.retval(d);
    const files = new Array<FileInfo>(d.filesLength());
    for (let i = 0; i < files.length; i += 1) {
      const f = d.files(i);
      files[i] = new FileInfo(<string>f.name(), <string>f.path(), f.isdir());
    }
    return new Directory(<string>d.name(), <string>d.path(), files);
  }
  case __std.FileSystemRetval.Error: {
    const err = new __std.Error();
    resp.retval(err);
    throw new Error(<string>err.message());
  }
  default:
    throw new Error('Unexpected response to fileinfo');
  }
}

export {
  info,
  dir,
};
