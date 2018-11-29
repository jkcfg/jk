import { __std as w } from '__std_Write_generated';
import { __std as m } from '__std_generated';
import flatbuffers from 'flatbuffers';

const Format = Object.freeze(w.Format);

function write(value, path = '', { format = Format.Auto, indent = 2 } = {}) {
  const builder = new flatbuffers.Builder(1024);
  const str = (format == Format.Raw) ? value.toString(): JSON.stringify(value);
  const strOff = builder.createString(str);
  const pathOff = builder.createString(path);

  w.WriteArgs.startWriteArgs(builder);
  w.WriteArgs.addValue(builder, strOff);
  w.WriteArgs.addPath(builder, pathOff);
  w.WriteArgs.addType(builder, format);
  w.WriteArgs.addIndent(builder, indent);
  const args = w.WriteArgs.endWriteArgs(builder);

  m.Message.startMessage(builder);
  m.Message.addArgsType(builder, m.Args.WriteArgs);
  m.Message.addArgs(builder, args);
  const message = m.Message.endMessage(builder);

  builder.finish(message);
  V8Worker2.send(builder.asArrayBuffer());
}

export {
  Format,
  write,
};
