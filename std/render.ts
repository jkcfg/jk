/**
 * @module: std/render
 */

import { requestAsPromise } from './internal/deferred';
import { flatbuffers } from './internal/flatbuffers';
import { __std } from './internal/__std_generated';

function uint8ToUint16Array(bytes: Uint8Array): Uint16Array {
  return new Uint16Array(bytes.buffer, bytes.byteOffset, bytes.byteLength / 2);
}

type Data = Uint8Array | string;
type Transform = (x: Data) => Data;

const compose = (f: Transform, g: Transform): Transform => (x: Data): Data => f(g(x));
const stringify = (bytes: Uint8Array): string => String.fromCodePoint(...uint8ToUint16Array(bytes));

export function render(pluginURL: string, params: object = {}): Promise<any> {
  const builder = new flatbuffers.Builder(512);
  const pluginURLOffset = builder.createString(pluginURL);
  const paramsStr = JSON.stringify(params);
  const paramsOffset = builder.createString(paramsStr);

  __std.RenderArgs.startRenderArgs(builder);
  __std.RenderArgs.addPluginURL(builder, pluginURLOffset);
  __std.RenderArgs.addParams(builder, paramsOffset);

  const argsOffset = __std.RenderArgs.endRenderArgs(builder);
  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.RenderArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const tx = compose(JSON.parse, stringify);
  return requestAsPromise((): null | ArrayBuffer => V8Worker2.send(builder.asArrayBuffer()), tx);
}
