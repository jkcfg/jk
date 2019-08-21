/**
 * @module std/internal
 */

import { __std } from './__std_generated';
import { flatbuffers } from './flatbuffers';
import { sendRequest, requestAsPromise } from './deferred';

// An RPC call
export function RPC(method: string, ...args: any[]): Promise<Uint8Array> {
  const builder = new flatbuffers.Builder(512);
  const argsOffsets = [];
  for (const arg of args) {
    let argType = __std.RPCValue.NONE;
    let argOffset = 0;
    if (arg instanceof Uint8Array) {
      let off = __std.RPCBytes.createBytesVector(builder, arg);
      __std.RPCBytes.startRPCBytes(builder);
      __std.RPCBytes.addBytes(builder, off);
      argOffset = __std.RPCBytes.endRPCBytes(builder);
      argType = __std.RPCValue.RPCBytes;
    } else {
      const serialisation = JSON.stringify(arg);
      let off = builder.createString(serialisation);
      __std.RPCSerialised.startRPCSerialised(builder);
      __std.RPCSerialised.addValue(builder, off);
      argOffset = __std.RPCSerialised.endRPCSerialised(builder);
      argType = __std.RPCValue.RPCSerialised;
    }
    __std.RPCArg.startRPCArg(builder);
    __std.RPCArg.addArgType(builder, argType);
    __std.RPCArg.addArg(builder, argOffset);
    argsOffsets.push(__std.RPCArg.endRPCArg(builder));
  }

  const methodOff = builder.createString(method);
  const argsOff = __std.RPCArgs.createArgsVector(builder, argsOffsets);
  __std.RPCArgs.startRPCArgs(builder);
  __std.RPCArgs.addMethod(builder, methodOff);
  __std.RPCArgs.addArgs(builder, argsOff);
  let off = __std.RPCArgs.endRPCArgs(builder);
  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.RPCArgs);
  __std.Message.addArgs(builder, off);
  off = __std.Message.endMessage(builder);
  builder.finish(off);
  return requestAsPromise(() => sendRequest(builder.asArrayBuffer()), (c: Uint8Array) => c);
}
