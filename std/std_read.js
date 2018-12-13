import { requestAsPromise } from 'std_deferred';
import flatbuffers from 'flatbuffers';
import { __std as r } from '__std_Read_generated';
import { __std as m } from '__std_generated';

const Encoding = Object.freeze(r.Encoding);

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
  r.ReadArgs.startReadArgs(builder);
  r.ReadArgs.addUrl(builder, urlOffset);
  let tx = bytes => bytes;
  switch (encoding) {
  case Encoding.UTF16:
    r.ReadArgs.addEncoding(builder, encoding);
    tx = stringify;
    break;
  case Encoding.JSON:
    r.ReadArgs.addEncoding(builder, encoding);
    tx = compose(JSON.parse, stringify);
    break;
  default:
  }
  const argsOffset = r.ReadArgs.endReadArgs(builder);

  m.Message.startMessage(builder);
  m.Message.addArgsType(builder, m.Args.ReadArgs);
  m.Message.addArgs(builder, argsOffset);
  const messageOffset = m.Message.endMessage(builder);
  builder.finish(messageOffset);
  return requestAsPromise(() => V8Worker2.send(builder.asArrayBuffer()), tx);
}

export {
  Encoding,
  read,
};
