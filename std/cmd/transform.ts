import * as std from '../index';
import * as param from '../param';
import { generate, File, GenerateParams } from './generate';
import { valuesFormatFromPath } from '../read';

type TransformFn = (value: any) => any | void;

const inputParams: GenerateParams = {
  stdout: param.Boolean('jk.transform.stdout', false),
  overwrite: param.Boolean('jk.transform.overwrite', false) ? std.Overwrite.Write : std.Overwrite.Err,
};

function transformOne(fn: TransformFn, file: string, obj: any): File {
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
    const format = valuesFormatFromPath(file);
    outputs.push(std.read(file, { format }).then((obj): File[] => {
      switch (format) {
      case std.Format.YAMLStream:
      case std.Format.JSONStream:
        return obj.map((v: any): File => transformOne(fn, file, v));
      default:
        return [transformOne(fn, file, obj)];
      }
    }));
  }
  generate(Promise.all(outputs).then((vs): File[] => Array.prototype.concat(...vs)), inputParams);
}

export default transform;
