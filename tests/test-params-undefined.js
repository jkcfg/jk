import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';

// Test that, if the user hasn't specified a parameter, we still get the
// default value (and not an error).
const p = param.String('nonexistent', 'success');
std.log(p);
