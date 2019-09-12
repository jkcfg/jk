import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';
import { generate } from '@jkcfg/std/generate';

const inputParams = {
  stdout: param.Boolean('jk.transform.stdout', false),
};

function readFormatFromPath(path) {
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

function transformOne(fn, file, obj) {
  let txObj = fn(obj);
  txObj = (txObj === undefined) ? obj : txObj;
  return {
    value: txObj,
    path: file,
  };
}

function transform(fn) {
  const inputFiles = param.Object('jk.transform.input', {});
  const outputs = [];
  for (const file of Object.keys(inputFiles)) {
    const format = readFormatFromPath(file);
    outputs.push(std.read(file, { format }).then((obj) => {
      switch (format) {
      case std.Format.YAMLStream:
      case std.Format.JSONStream:
        return obj.map(v => transformOne(fn, file, v));
      default:
        return [transformOne(fn, file, obj)];
      }
    }));
  }
  generate(Promise.all(outputs).then(vs => Array.prototype.concat(...vs)), inputParams);
}

export default transform;
