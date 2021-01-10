import { Format, Overwrite, read, stdin, print } from '../index';
import * as host from '@jkcfg/std/internal/host'; // magic module
import * as param from '../param';
import { generate, File, GenerateParams, maybeSetFormat } from './generate';
import { valuesFormatFromPath, valuesFormatFromExtension } from '../read';

type TransformFn = (value: any) => any | void;

const generateParams: GenerateParams = {
  stdout: param.Boolean('jk.transform.stdout', false),
  overwrite: param.Boolean('jk.transform.overwrite', false) ? Overwrite.Write : Overwrite.Err,
};
maybeSetFormat(generateParams, param.String('jk.generate.format', undefined)); // NB jk.generate. param

// If we're told to overwrite, we need to be able to write to the
// files mentioned on the command-line; but not otherwise.
if (generateParams.overwrite == Overwrite.Write) {
  generateParams.writeFile = host.write;
}

function transform(fn: TransformFn): void {

  function transformOne(obj: any): any {
    let txObj = fn(obj);
    txObj = (txObj === undefined) ? obj : txObj;
    return txObj;
  }

  const inputFiles = param.Object('jk.transform.input', {});
  const outputs = [];

  for (const path of Object.keys(inputFiles)) {
    if (path === '') { // read from stdin
      const stdinFormat = param.String('jk.transform.stdin.format', 'yaml');
      const format = valuesFormatFromExtension(stdinFormat);
      const path = `stdin.${stdinFormat}`;  // path is a stand-in
      const value = read(stdin, { format }).then(v => v.map(transformOne));
      outputs.push({ path, value, format });
      continue;
    }

    const format = valuesFormatFromPath(path);
    outputs.push(host.read(path, { format }).then((obj): File => {
      switch (format) {
      case Format.YAMLStream:
      case Format.JSONStream:
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
  generate(Promise.all(outputs), generateParams);
}

export default transform;
