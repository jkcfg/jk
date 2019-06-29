import { flatbuffers } from './internal/flatbuffers';
import { __std } from './internal/__std_generated';

import ParamType = __std.ParamType;

function getParameter<T>(type: ParamType, path: string, defaultValue: T): T | undefined {
  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
  const isObject = type === __std.ParamType.Object;
  const defaultValueOffset = isObject
    && defaultValue !== undefined
    && builder.createString(JSON.stringify(defaultValue));

  __std.ParamArgs.startParamArgs(builder);
  __std.ParamArgs.addPath(builder, pathOffset);
  __std.ParamArgs.addType(builder, type);
  if (isObject) {
    __std.ParamArgs.addDefaultValue(builder, defaultValueOffset);
  }
  const argsOffset = __std.ParamArgs.endParamArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.ParamArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);

  const bytes = <ArrayBuffer>V8Worker2.send(builder.asArrayBuffer());

  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = __std.ParamResponse.getRootAsParamResponse(buf);

  switch (resp.retvalType()) {
  case __std.ParamRetval.ParamValue: {
    // The runtime has returned the parameter value (or null).
    const ret = new __std.ParamValue();
    resp.retval(ret);

    const v = JSON.parse(<string>ret.json());
    if (v == null) {
      return defaultValue;
    }
    return v;
  }
  case __std.ParamRetval.Error: {
    // The runtime has returned an error.
    const err = new __std.Error();
    resp.retval(err);
    throw new Error(<string>err.message());
  }
  default:
    throw new Error('Unexpected response to param');
  }
}

export function Boolean(path: string, defaultValue?: boolean): boolean | undefined {
  return getParameter(__std.ParamType.Boolean, path, defaultValue);
}

export function Number(path: string, defaultValue?: number): number | undefined {
  return getParameter(__std.ParamType.Number, path, defaultValue);
}

export function String(path: string, defaultValue?: string): string | undefined {
  return getParameter(__std.ParamType.String, path, defaultValue);
}

export function Object(path: string, defaultValue?: object): object | undefined {
  return getParameter(__std.ParamType.Object, path, defaultValue);
}

export function all(): object {
  return <object>Object('');
}

export const param = {
  all,
  Boolean,
  Number,
  String,
  Object,
};
