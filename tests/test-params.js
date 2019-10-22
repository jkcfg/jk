import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';

const b = param.Boolean('myBoolean', false);
const n = param.Number('myNumber', 3.14);
const s = param.String('myString', 'foo');

std.print({
  myBoolean: b,
  myNumber: n,
  myString: s,
});

// Test default values are shining through when the parameters aren't
// specified.
const bD = param.Boolean('myBooleanD', false);
const nD = param.Number('myNumberD', 3.14);
const sD = param.String('myStringD', 'foo');
const oD = param.Object('myObjectD', {
  s: 'bar',
  b: true,
  o: {
    xxx: 'yyy',
  },
});

std.print({
  myBoolean: bD,
  myNumber: nD,
  myString: sD,
  myObject: oD,
});
