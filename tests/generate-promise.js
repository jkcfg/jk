import * as std from '@jkcfg/std';

function f(filename) {
  return std.read(filename).then(o => [
    { path: 'object.yaml', value: o },
  ]);
}

export default f('success.json');
