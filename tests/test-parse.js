import { parse, log, Format } from '@jkcfg/std';

const foo = parse('{ "foo" : 1 }', Format.JSON);
log(foo);
