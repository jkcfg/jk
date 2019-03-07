import std from '@jkcfg/std';
import resource1 from './test-module-resource/resource';
import resource2 from './test-module-resource/submodule/resource';
import contents from './test-module-resource/fs';

Promise.all([resource1, resource2, contents]).then(resources => resources.forEach(std.log));
