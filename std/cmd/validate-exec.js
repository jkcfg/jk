import validate from '@jkcfg/std/cmd/validate';

const makeValidateFn = new Function(`
  return (%s);
`);

const fn = makeValidateFn();
validate(fn);
