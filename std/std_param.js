import flatbuffers from 'flatbuffers';
import { __std } from '__std_generated';

function getParameter(type, path, defaultValue) {
  const builder = new flatbuffers.Builder(512);
  const pathOffset = builder.createString(path);
  const isObject = type === __std.ParamType.Object;
  const defaultValueOffset = isObject && builder.createString(JSON.stringify(defaultValue));

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

  const bytes = V8Worker2.send(builder.asArrayBuffer());

  const buf = new flatbuffers.ByteBuffer(new Uint8Array(bytes));
  const resp = __std.ParamResponse.getRootAsParamResponse(buf);
  const v = JSON.parse(resp.json());
  if (v == null) {
    return defaultValue;
  }
  return v;
}

function getBoolean(path, defaultValue) {
  return getParameter(__std.ParamType.Boolean, path, defaultValue);
}

function getNumber(path, defaultValue) {
  return getParameter(__std.ParamType.Number, path, defaultValue);
}

function getString(path, defaultValue) {
  return getParameter(__std.ParamType.String, path, defaultValue);
}

function getObject(path, defaultValue) {
  return getParameter(__std.ParamType.Object, path, defaultValue);
}

export const param = {
  Boolean: getBoolean,
  Number: getNumber,
  String: getString,
  Object: getObject,
};
