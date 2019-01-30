// This test checks that a read path is resolved relative to an input
// directory. It's used from more than one test -- other tests may be
// given as .cmd files invoking this file.

import std from '@jkcfg/std';

// the location of this file will differ depending on the
// `--input-directory` given (or not given) to jk run
std.read('input.json').then(std.log)
