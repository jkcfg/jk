import { requestAsPromise } from 'std_deferred';
import flatbuffers from 'flatbuffers';
import { __std } from '__std_generated';
import { Format } from 'std_write';

const Encoding = Object.freeze(__std.Encoding);

function uint8ToUint16Array(bytes) {
  return new Uint16Array(bytes.buffer, bytes.byteOffset, bytes.byteLength / 2);
}

const compose = (f, g) => x => f(g(x));
const stringify = bytes => String.fromCodePoint(...uint8ToUint16Array(bytes));

// read requests the path and returns a promise that will be resolved
// with the contents at the path, or rejected.
function read(path, { encoding = Encoding.UTF16, format = Format.Auto } = {}) {
  const builder = new flatbuffers.Builder(512);
  const urlOffset = builder.createString(path);
  __std.ReadArgs.startReadArgs(builder);
  __std.ReadArgs.addUrl(builder, urlOffset);
  __std.ReadArgs.addEncoding(builder, encoding);
  __std.ReadArgs.addFormat(builder, format);
  const argsOffset = __std.ReadArgs.endReadArgs(builder);
  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.ReadArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  let tx = bytes => bytes;
  switch (encoding) {
  case Encoding.UTF16:
    tx = stringify;
    break;
  case Encoding.JSON:
    tx = compose(JSON.parse, stringify);
    break;
  default:
    break;
  }

  return requestAsPromise(() => V8Worker2.send(builder.asArrayBuffer()), tx);
}

export {
  Encoding,
  read,
};
