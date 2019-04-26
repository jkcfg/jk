import { requestAsPromise } from './deferred';
import { flatbuffers } from './flatbuffers';
import { __std } from './__std_generated';

import Encoding = __std.Encoding;
import Format = __std.Format;

function uint8ToUint16Array(bytes: Uint8Array): Uint16Array {
  return new Uint16Array(bytes.buffer, bytes.byteOffset, bytes.byteLength / 2);
}

type Data = Uint8Array | string;
type Transform = (x: Data) => Data;

const compose = (f: Transform, g: Transform): Transform => (x: Data): Data => f(g(x));
const stringify = (bytes: Uint8Array): string => String.fromCodePoint(...uint8ToUint16Array(bytes));

interface ReadOptions {
  encoding?: Encoding;
  format?: Format;
  module?: string;
}

// read requests the path and returns a promise that will be resolved
// with the contents at the path, or rejected.
function read(path: string, opts: ReadOptions = {}): Promise<any> {
  const { encoding = Encoding.JSON, format = Format.FromExtension, module } = opts;

  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
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

  let tx: Transform = (bytes: Uint8Array): Uint8Array => bytes;
  switch (encoding) {
  case Encoding.String:
    tx = stringify;
    break;
  case Encoding.JSON:
    tx = compose(JSON.parse, stringify);
    break;
  default:
    break;
  }

  return requestAsPromise((): null | ArrayBuffer => V8Worker2.send(builder.asArrayBuffer()), tx);
}

export {
  Encoding,
  read,
};
