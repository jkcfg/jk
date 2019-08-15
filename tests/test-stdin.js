import * as std from '@jkcfg/std';

std.read('', { format: std.Format.JSONStream })
  .then(v => std.log(v, { format: std.Format.YAMLStream }));
