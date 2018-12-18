import { __std } from '__std_generated';
import flatbuffers from 'flatbuffers';

const Format = Object.freeze(__std.Format);

function write(value, path = '', { format = Format.Auto, indent = 2, override = true } = {}) {
  const builder = new flatbuffers.Builder(1024);
  const str = (format === Format.Raw) ? value.toString() : JSON.stringify(value);
  const strOff = builder.createString(str);
  const pathOff = builder.createString(path);

  __std.WriteArgs.startWriteArgs(builder);
  __std.WriteArgs.addValue(builder, strOff);
  __std.WriteArgs.addPath(builder, pathOff);
  __std.WriteArgs.addFormat(builder, format);
  __std.WriteArgs.addIndent(builder, indent);
  __std.WriteArgs.addOverride(builder, override);
  const args = __std.WriteArgs.endWriteArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.WriteArgs);
  __std.Message.addArgs(builder, args);
  const message = __std.Message.endMessage(builder);

  builder.finish(message);
  V8Worker2.send(builder.asArrayBuffer());
}

export {
  Format,
  write,
};
