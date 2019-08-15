import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';

const b = param.Boolean('myBoolean', false);
const n = param.Number('myNumber', 3.14);
const s = param.String('myString', 'foo');
const o = param.Object('myObject', {
  s: 'bar',
  b: true,
  o: {
    xxx: 'yyy',
  },
});

std.log({
  myBoolean: b,
  myNumber: n,
  myString: s,
  myObject: o,
});
