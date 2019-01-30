/* eslint "import/no-unresolved": [2, { ignore: ['.*'] }] */
// ^ switch this rule in eslint off, since module resolution is what's under test here.
/* eslint "import/first": [0] */
/* eslint "import/newline-after-import": [0] */
import std from '@jkcfg/std';

// node_modules/testcase/index.js
import msg1 from 'testcase';
std.write(msg1, 'test1.json');

// node_modules/testcase/submodule.js
import msg2 from 'testcase/submodule';
std.write(msg2, 'test2.json');

// node_modules/testcase/indirect.js imports ./test3
import msg3 from 'testcase/indirect';
std.write(msg3, 'test3.json');

// node_modules/testcase/subdir/package.json specifies test4.js in its
// `module` field
import msg4 from 'testcase/subdir';
std.write(msg4, 'test4.json');

// node_modules/testcase/vendor.js imports from 'vendor', which is to
// be found at node_modules/testcase/node_modules/vendored/index.js
import msg5 from 'testcase/vendor';
std.write(msg5, 'test5.json');
