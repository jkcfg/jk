import * as std from '../index';
import * as param from '../param';
import { generate, File, GenerateParams } from './generate';
import { valuesFormatFromPath } from '../read';

type TransformFn = (value: any) => any | void;

const inputParams: GenerateParams = {
  stdout: param.Boolean('jk.transform.stdout', false),
  overwrite: param.Boolean('jk.transform.overwrite', false) ? std.Overwrite.Write : std.Overwrite.Err,
};

function transform(fn: TransformFn): void {

  function transformOne(obj: any): any {
    let txObj = fn(obj);
    txObj = (txObj === undefined) ? obj : txObj;
    return txObj;
  }

  const inputFiles = param.Object('jk.transform.input', {});
  const outputs = [];
  for (const path of Object.keys(inputFiles)) {
    const format = valuesFormatFromPath(path);
    outputs.push(std.read(path, { format }).then((obj): File => {
      switch (format) {
      case std.Format.YAMLStream:
      case std.Format.JSONStream:
        return {
          path,
          format,
          value: Array.prototype.map.call(obj, transformOne),
        };
      default:
        return { path, value: transformOne(obj) };
      }
    }));
  }
  generate(Promise.all(outputs), inputParams);
}

export default transform;
