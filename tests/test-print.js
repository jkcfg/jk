import { print, Format } from '@jkcfg/std';
import doTest from './print-and-log';

doTest(print);

// And tune the output indentation
print({ kind: 'Bar', foo: 1.2 }, { indent: 4 });

// Key order is deterministic.
print({ kind: 'Bar', foo: 1.2 });
print({ foo: 1.2, kind: 'Bar' });

// Objects, but in YAML
print({ kind: 'Bar', foo: { number: 1.2, string: 'mystring' } }, { format: Format.YAML });
print({ foo: { number: 1.2, string: 'mystring' }, kind: 'Bar' }, { format: Format.YAML });
