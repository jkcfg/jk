import std from '@jkcfg/std';

import msg1 from 'testcase'; // node_modules/testcase/index.js
std.write(msg1, 'test1.json');

import msg2 from 'testcase/submodule'; // node_modules/testcase/submodule.js
std.write(msg2, 'test2.json');

import msg3 from 'testcase/indirect'; // node_modules/testcase/indirect.js imports ./test3
std.write(msg3, 'test3.json');

import msg4 from 'testcase/subdir'; // node_modules/testcase/subdir/package.json specifies test4.js in its `module` field
std.write(msg4, 'test4.json');
