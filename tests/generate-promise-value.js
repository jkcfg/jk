import * as std from '@jkcfg/std';

function f(filename) {
  return std.read(filename).then(o => o);
}

export default [
  { path: 'object.yaml', value: f('success.json') },
];
