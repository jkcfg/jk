import std from '@jkcfg/std';

import msg1 from 'testcase'; // node_modules/testcase/index.js
std.write(msg1, 'test1.json');

import msg2 from 'testcase/submodule'; // node_modules/testcase/submodule.js
std.write(msg2, 'test2.json');

