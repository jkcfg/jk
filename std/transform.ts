import * as std from './index';
import * as param from './param';
import { generate, Value, GenerateParams } from './generate';

type TransformFn = (value: any) => any | void;

const inputParams: GenerateParams = {
  stdout: param.Boolean('jk.transform.stdout', false),
  overwrite: param.Boolean('jk.transform.overwrite', false) ? std.Overwrite.Write : std.Overwrite.Err,
};

function readFormatFromPath(path: string): std.Format {
  const ext = path.split('.').pop();
  switch (ext) {
  case 'yaml':
  case 'yml':
    return std.Format.YAMLStream;
  case 'json':
    return std.Format.JSONStream;
  default:
    return std.Format.FromExtension;
  }
}

function transformOne(fn: TransformFn, file: string, obj: any): Value {
  let txObj = fn(obj);
  txObj = (txObj === undefined) ? obj : txObj;
  return {
    value: txObj,
    path: file,
  };
}

function transform(fn: TransformFn): void {
  const inputFiles = param.Object('jk.transform.input', {});
  const outputs = [];
  for (const file of Object.keys(inputFiles)) {
    const format = readFormatFromPath(file);
    outputs.push(std.read(file, { format }).then((obj): Value[] => {
      switch (format) {
      case std.Format.YAMLStream:
      case std.Format.JSONStream:
        return obj.map((v: any): Value => transformOne(fn, file, v));
      default:
        return [transformOne(fn, file, obj)];
      }
    }));
  }
  generate(Promise.all(outputs).then((vs): Value[] => Array.prototype.concat(...vs)), inputParams);
}

export default transform;
