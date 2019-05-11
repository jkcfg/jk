import * as std from '@jkcfg/std';

function f(filename) {
  return std.read(filename).then(o => [
    { file: 'object.yaml', content: o },
  ]);
}

export default f('success.json');
