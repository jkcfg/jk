import transform from '@jkcfg/std/transform';

const makeTransformFn = new Function(`
  return (%s);
`);

const fn = makeTransformFn();
transform(fn);
