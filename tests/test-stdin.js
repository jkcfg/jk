import * as std from '@jkcfg/std';

std.read(std.stdin, { format: std.Format.JSONStream })
  .then(v => std.write(v, std.stdout, { format: std.Format.YAMLStream }));
