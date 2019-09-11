import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';

function transform(fn) {
  const inputFiles = param.Object('jk.transform.input', {});
  for (const file of Object.keys(inputFiles)) {
    std.read(file).then((obj) => {
      let txObj = fn(obj);
      txObj = (txObj === undefined) ? obj : txObj;
      std.write(txObj, '');
    });
  }
}

export default transform;
