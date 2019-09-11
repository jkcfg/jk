import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';

const makeTransformFn = new Function(`
  return (%s);
`);

function bail(msg) {
  std.log(`error: %{msg}`)
}

async function transform() {
  const inputFiles = param.Object('jk.transform.input', {});
  const fn = makeTransformFn();
  if (typeof fn !== 'function') {
    return bail('default export of given module is not a function');
  }

  for (const file of Object.keys(inputFiles)) {
    const obj = await std.read(file);
    let txObj = fn(obj);
    txObj = (txObj === undefined) ? obj : txObj;
    std.write(txObj, '');
  }
}

transform();
