import { __std as w } from '__std_Write_generated';
import { __std as m } from '__std_generated';
import flatbuffers from 'flatbuffers';

const outputFormat = w.OutputFormat;

function write(value, path = '', fmt = outputFormat.Auto) {
  const builder = new flatbuffers.Builder(1024);
  const json = JSON.stringify(value);
  const jsonStr = builder.createString(json);
  const pathStr = builder.createString(path);

  w.WriteArgs.startWriteArgs(builder);
  w.WriteArgs.addValue(builder, jsonStr);
  w.WriteArgs.addPath(builder, pathStr);
  w.WriteArgs.addType(builder, fmt);
  const args = w.WriteArgs.endWriteArgs(builder);

  m.Message.startMessage(builder);
  m.Message.addArgsType(builder, m.Args.WriteArgs);
  m.Message.addArgs(builder, args);
  const message = m.Message.endMessage(builder);

  builder.finish(message);
  V8Worker2.send(builder.asArrayBuffer());
}

export {
  outputFormat,
  write,
};
