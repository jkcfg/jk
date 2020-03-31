/**
 * @module std
 */

import { __std } from './internal/__std_generated';
import { flatbuffers } from './internal/flatbuffers';
import { sendRequest } from './internal/deferred';

export const stdout: unique symbol = Symbol('<stdout>');

/* we re-define Format from the generated __std.Format to document it */

export enum Format {
  FromExtension= 0,
  JSON= 1,
  YAML= 2,
  Raw= 3,
  YAMLStream= 4,
  JSONStream= 5,
  HCL= 6,
}

export enum Overwrite {
  Skip= 0,
  Write= 1,
  Err= 2,
}

export interface WriteOptions {
  format?: Format;
  indent?: number;
  overwrite?: Overwrite | boolean;
  module?: string;
}

type WritePath = string | typeof stdout;

export function write(value: any, path: WritePath = stdout, opts: WriteOptions = {}): void {
  if (value === undefined) {
    throw TypeError('cannot write undefined value');
  }

  const {
    format = Format.FromExtension,
    indent = 2,
    overwrite = Overwrite.Write,
    module,
  } = opts;
  const pathArg = (path === stdout) ? '' : path;

  let overwriteVal: Overwrite;
  if (typeof overwrite === 'boolean') {
    overwriteVal = overwrite ? Overwrite.Write : Overwrite.Skip;
  } else {
    overwriteVal = overwrite;
  }

  const builder = new flatbuffers.Builder(1024);
  const str = (format === Format.Raw) ? value.toString() : JSON.stringify(value);
  const strOffset = builder.createString(str);
  const pathOffset = builder.createString(pathArg);
  let moduleOffset = 0;
  if (module !== undefined) {
    moduleOffset = builder.createString(module);
  }

  __std.WriteArgs.startWriteArgs(builder);
  __std.WriteArgs.addValue(builder, strOffset);
  __std.WriteArgs.addPath(builder, pathOffset);
  __std.WriteArgs.addFormat(builder, format);
  __std.WriteArgs.addIndent(builder, indent);
  __std.WriteArgs.addOverwrite(builder, overwriteVal);
  if (module !== undefined) {
    __std.WriteArgs.addModule(builder, moduleOffset);
  }
  const args = __std.WriteArgs.endWriteArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.WriteArgs);
  __std.Message.addArgs(builder, args);
  const message = __std.Message.endMessage(builder);

  builder.finish(message);
  const buf = sendRequest(builder.asArrayBuffer());
  if (buf === undefined) {
    return;
  }
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const resp = __std.Error.getRootAsError(data);
  throw new Error(resp.message());
}

// print is a convenience for printing any value to stdout
export function print(value: any, opts: WriteOptions = {}): void {
  if (arguments.length === 0) {
    write('\n', stdout, { format: Format.Raw });
    return;
  }
  if (value === undefined) {
    write('undefined\n', stdout, { format: Format.Raw });
    return;
  }
  write(value, stdout, opts);
}
