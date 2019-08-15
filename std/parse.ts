/**
 * @module std
 */

import { flatbuffers } from './internal/flatbuffers';
import { __std } from './internal/__std_generated';

import Format = __std.Format;

export function parse(input: string, format?: Format): any {
  const builder = new flatbuffers.Builder(512);
  const inputOffset = builder.createString(input);

  __std.ParseArgs.startParseArgs(builder);
  __std.ParseArgs.addChars(builder, inputOffset);
  if (format !== undefined) {
    __std.ParseArgs.addFormat(builder, format);
  }
  const argsOffset = __std.ParseArgs.endParseArgs(builder);
  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.ParseArgs);
  __std.Message.addArgs(builder, argsOffset);
  builder.finish(__std.Message.endMessage(builder));

  const buf = V8Worker2.send(builder.asArrayBuffer());
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const resp = __std.ParseUnparseResponse.getRootAsParseUnparseResponse(data);
  switch (resp.retvalType()) {
  case __std.ParseUnparseRetval.Error: {
    const err = new __std.Error();
    resp.retval(err);
    throw new Error(err.message());
  }
  case __std.ParseUnparseRetval.ParseUnparseData: {
    const val = new __std.ParseUnparseData();
    resp.retval(val);
    return JSON.parse(val.data());
  }
  default:
    throw new Error('Response type was not set');
  }
}

export function unparse(obj: any, format?: Format): string {
  const builder = new flatbuffers.Builder(512);
  const inputOffset = builder.createString(JSON.stringify(obj));

  __std.UnparseArgs.startUnparseArgs(builder);
  __std.UnparseArgs.addObject(builder, inputOffset);
  if (format !== undefined) {
    __std.UnparseArgs.addFormat(builder, format);
  }
  const argsOffset = __std.UnparseArgs.endUnparseArgs(builder);
  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.UnparseArgs);
  __std.Message.addArgs(builder, argsOffset);
  builder.finish(__std.Message.endMessage(builder));

  const buf = V8Worker2.send(builder.asArrayBuffer());
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const resp = __std.ParseUnparseResponse.getRootAsParseUnparseResponse(data);
  switch (resp.retvalType()) {
  case __std.ParseUnparseRetval.Error: {
    const err = new __std.Error();
    resp.retval(err);
    throw new Error(err.message());
  }
  case __std.ParseUnparseRetval.ParseUnparseData: {
    const val = new __std.ParseUnparseData();
    resp.retval(val);
    return val.data();
  }
  default:
    throw new Error('Response type was not set');
  }
}
