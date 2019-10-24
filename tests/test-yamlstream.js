import * as std from '@jkcfg/std';

const value = std.read('./test-yamlstream.js.expected', { format: std.Format.YAMLStream });
value.then(v => std.print(v, { format: std.Format.YAMLStream }));
