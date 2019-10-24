import transform from '@jkcfg/std/cmd/transform';

const makeTransformFn = new Function(`
  return (%s);
`);

const fn = makeTransformFn();
transform(fn);
