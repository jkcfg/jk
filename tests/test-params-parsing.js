import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';

for (const name of ['foo.1', 'foo.2', 'foo.3', 'foo.4', 'foo.5', 'foo.6']) {
  std.log(`${name}: ${param.String(name)}`);
}
