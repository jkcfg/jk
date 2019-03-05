import std from '@jkcfg/std';
import resource1 from './test-module-resource/resource';
import resource2 from './test-module-resource/submodule/resource';

Promise.all([resource1, resource2]).then(resources => resources.forEach(std.log));
