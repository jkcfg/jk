import * as rxjs from 'node_modules/rxjs/_esm2015/index.js'
import write from 'write.js'

rxjs.of("foo", "bar", "baz").forEach(write);
