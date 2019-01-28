/* eslint "import/no-unresolved": [2, { ignore: ['.*'] }] */
// ^ switch this rule in eslint off, since module resolution is what's under test here.
/* eslint "import/first": [0] */
/* eslint "import/newline-after-import": [0] */
import std from '@jkcfg/std';

import msg1 from 'testcase'; // node_modules/testcase/index.js
std.write(msg1, 'test1.json');

import msg2 from 'testcase/submodule'; // node_modules/testcase/submodule.js
std.write(msg2, 'test2.json');

import msg3 from 'testcase/indirect'; // node_modules/testcase/indirect.js imports ./test3
std.write(msg3, 'test3.json');

import msg4 from 'testcase/subdir'; // node_modules/testcase/subdir/package.json specifies test4.js in its `module` field
std.write(msg4, 'test4.json');

import msg5 from 'testcase/vendor';
// node_modules/testcase/vendor.js imports from 'vendor', which is to be found at node_modules/testcase/node_modules/vendored/index.js
std.write(msg5, 'test5.json');
