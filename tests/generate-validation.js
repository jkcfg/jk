import * as param from '@jkcfg/std/param';

function fooIsBar(v) {
  if (v.foo === 'bar') {
    return 'ok';
  }
  return ['foo is not bar'];
}

export default [
  { path: 'foo.yaml', value: { foo: param.String('foo', 'bar') }, validate: fooIsBar },
];
