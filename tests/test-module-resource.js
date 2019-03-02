import std from '@jkcfg/std';
import resource1 from './test-module-resource/resource';
import resource2 from './test-module-resource/submodule/resource';

resource1.then(std.log);
resource2.then(std.log);
