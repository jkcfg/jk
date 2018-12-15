import { requestAsPromise } from 'std_deferred';
import flatbuffers from 'flatbuffers';
import { __std } from '__std_generated';

const Encoding = Object.freeze(__std.Encoding);

function uint8ToUint16Array(bytes) {
  return new Uint16Array(bytes.buffer, bytes.byteOffset, bytes.byteLength / 2);
}

const compose = (f, g) => x => f(g(x));
const stringify = bytes => String.fromCodePoint(...uint8ToUint16Array(bytes));

// read requests a URL and returns a promise that will be resolved
// with the contents at the URL, or rejected.
function read(url, { encoding = Encoding.Bytes } = {}) {
  const builder = new flatbuffers.Builder(512);
  const urlOffset = builder.createString(url);
  __std.ReadArgs.startReadArgs(builder);
  __std.ReadArgs.addUrl(builder, urlOffset);
  let tx = bytes => bytes;
  switch (encoding) {
  case Encoding.UTF16:
    __std.ReadArgs.addEncoding(builder, encoding);
    tx = stringify;
    break;
  case Encoding.JSON:
    __std.ReadArgs.addEncoding(builder, encoding);
    tx = compose(JSON.parse, stringify);
    break;
  default:
  }
  const argsOffset = __std.ReadArgs.endReadArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.ReadArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);
  return requestAsPromise(() => V8Worker2.send(builder.asArrayBuffer()), tx);
}

export {
  Encoding,
  read,
};
