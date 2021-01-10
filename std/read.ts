/**
 * @module std
 */

import { requestAsPromise, sendRequest } from './internal/deferred';
import {
  Transform,
  ident,
  stringFromUTF16Bytes,
  valueFromUTF16Bytes,
} from './internal/data';
import { flatbuffers } from './internal/flatbuffers';
import { __std } from './internal/__std_generated';
import { Format } from './write';

export const stdin: unique symbol = Symbol('<stdin>');

/* we re-define Encoding from the generated __std.Encoding to document it */

export enum Encoding {
  Bytes= 0,
  String= 1,
  JSON= 2,
}

export interface ReadOptions {
  encoding?: Encoding;
  format?: Format;
  module?: string;
}

// splitPath returns [all-but-extension, extension] for a path. If a
// path does not end with an extension, it will be an empty string.
export function splitPath(path: string): [string, string] {
  const parts = path.split('.');
  const ext = parts.pop();
  // When there's no extension, either there will be a single part (no
  // dots anywhere), or a path separator in the last part (a dot
  // somewhere before the last path segment)
  if (parts.length === 0 || ext.includes('/')) {
    return [ext, ''];
  }
  return [parts.join(''), ext];
}

function extension(path: string): string {
  return splitPath(path)[1];
}

// formatFromPath guesses, for a file path, the format in which to
// read the file. It will assume one value per file, so if you have
// files that may have multiple values (e.g., YAML streams), it's
// better to use `valuesFormatFromPath` and be prepared to get
// multiple values.
export function formatFromPath(path: string): Format {
  switch (extension(path)) {
  case 'yaml':
  case 'yml':
    return Format.YAML;
  case 'json':
    return Format.JSON;
  case 'hcl':
  case 'tf':
    return Format.HCL;
  default:
    return Format.JSON;
  }
}

// valuesFormatFromExtension returns the format implied by a
// particular file extension.
export function valuesFormatFromExtension(ext: string): Format {
  switch (ext) {
  case 'yaml':
  case 'yml':
    return Format.YAMLStream;
  case 'json':
    return Format.JSONStream;
  default:
    return Format.FromExtension;
  }
}

// valuesFormatFromPath guesses, for a path, the format that will
// return all values in a file. In other words, it prefers YAML
// streams and concatenated JSON. You may need to treat the read value
// differently depending on the format you got here, since YAMLStream
// and JSONStream will both result in an array of values.
export function valuesFormatFromPath(path: string): Format {
  const ext = extension(path);
  return valuesFormatFromExtension(ext);
}

type ReadPath = string | typeof stdin;

// read requests the path and returns a promise that will be resolved
// with the contents at the path, or rejected.
export function read(path: ReadPath = stdin, opts: ReadOptions = {}): Promise<any> {
  const { encoding = Encoding.JSON, format = Format.FromExtension, module } = opts;
  const pathArg = (path === stdin) ? '' : path;

  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(pathArg);
  let moduleOffset = 0;
  if (module !== undefined) {
    moduleOffset = builder.createString(module);
  }
  __std.ReadArgs.startReadArgs(builder);
  __std.ReadArgs.addPath(builder, pathOffset);
  __std.ReadArgs.addEncoding(builder, encoding);
  __std.ReadArgs.addFormat(builder, format);
  if (module !== undefined) {
    __std.ReadArgs.addModule(builder, moduleOffset);
  }
  const argsOffset = __std.ReadArgs.endReadArgs(builder);
  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.ReadArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  let tx: Transform = ident;
  switch (encoding) {
  case Encoding.String:
    tx = stringFromUTF16Bytes;
    break;
  case Encoding.JSON:
    tx = valueFromUTF16Bytes;
    break;
  default:
    break;
  }

  return requestAsPromise((): null | ArrayBuffer => sendRequest(builder.asArrayBuffer()), tx);
}
