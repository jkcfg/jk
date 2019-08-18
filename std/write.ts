/**
 * @module std
 */

import { __std } from './internal/__std_generated';
import { flatbuffers } from './internal/flatbuffers';
import { sendRequest } from './internal/deferred';

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

export interface WriteOptions {
  format?: Format;
  indent?: number;
  overwrite?: boolean;
}

export function write(value: any, path = '', { format = Format.FromExtension, indent = 2, overwrite = true }: WriteOptions = {}): void {
  if (value === undefined) {
    throw TypeError('cannot write undefined value');
  }

  const builder = new flatbuffers.Builder(1024);
  const str = (format === Format.Raw) ? value.toString() : JSON.stringify(value);
  const strOff = builder.createString(str);
  const pathOff = builder.createString(path);

  __std.WriteArgs.startWriteArgs(builder);
  __std.WriteArgs.addValue(builder, strOff);
  __std.WriteArgs.addPath(builder, pathOff);
  __std.WriteArgs.addFormat(builder, format);
  __std.WriteArgs.addIndent(builder, indent);
  __std.WriteArgs.addOverwrite(builder, overwrite);
  const args = __std.WriteArgs.endWriteArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.WriteArgs);
  __std.Message.addArgs(builder, args);
  const message = __std.Message.endMessage(builder);

  builder.finish(message);
  sendRequest(builder.asArrayBuffer());
}
