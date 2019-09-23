import { write, Overwrite } from '@jkcfg/std';

write({ foo: 1 }, 'test-overwrite.js', { overwrite: Overwrite.Err });
