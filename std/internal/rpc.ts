/**
 * @module std/internal
 */

import { __std } from './__std_generated';
import { flatbuffers } from './flatbuffers';
import { sendRequest, requestAsPromise } from './deferred';
import { ident } from './data';

function encode(method: string, args: any[], sync: boolean): ArrayBuffer {
  const builder = new flatbuffers.Builder(512);
  const argsOffsets = [];
  for (const arg of args) {
    let argType = __std.RPCValue.NONE;
    let argOffset = 0;
    if (arg instanceof Uint8Array) {
      let off = __std.Data.createBytesVector(builder, arg);
      __std.Data.startData(builder);
      __std.Data.addBytes(builder, off);
      argOffset = __std.Data.endData(builder);
      argType = __std.RPCValue.Data;
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
  __std.RPCArgs.addSync(builder, sync);
  let off = __std.RPCArgs.endRPCArgs(builder);
  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.RPCArgs);
  __std.Message.addArgs(builder, off);
  off = __std.Message.endMessage(builder);
  builder.finish(off);
  return builder.asArrayBuffer();
}

// An asynchronous RPC call
export function RPC(method: string, ...args: any[]): Promise<Uint8Array> {
  return requestAsPromise(() => sendRequest(encode(method, args, false)), ident);
}

// A synchronous RPC call
export function RPCSync(method: string, ...args: any[]): Uint8Array {
  const result = sendRequest(encode(method, args, true));
  const buffer = new flatbuffers.ByteBuffer(new Uint8Array(result));
  const resp = __std.RPCSyncResponse.getRootAsRPCSyncResponse(buffer);
  switch (resp.retvalType()) {
  case __std.RPCSyncRetval.Error: {
    const err = new __std.Error();
    resp.retval(err);
    throw new Error(err.message());
  }
  case __std.RPCSyncRetval.Data: {
    const data = new __std.Data();
    resp.retval(data);
    return data.bytesArray();
  }
  default:
    throw new Error('unable to decode response');
  }
}
