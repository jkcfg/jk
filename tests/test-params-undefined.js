import std from '@jkcfg/std';

// Test that, if the user hasn't specified a parameter, we still get the
// default value (and not an error).
const p = std.param.String('nonexistent', 'success');
std.log(p);
