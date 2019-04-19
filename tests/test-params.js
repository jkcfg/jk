import std from '@jkcfg/std';

const b = std.param.Boolean('myBoolean', false);
const n = std.param.Number('myNumber', 3.14);
const s = std.param.String('myString', 'foo');

std.log({
  myBoolean: b,
  myNumber: n,
  myString: s,
});

// Test default values are shining through when the parameters aren't
// specified.
const bD = std.param.Boolean('myBooleanD', false);
const nD = std.param.Number('myNumberD', 3.14);
const sD = std.param.String('myStringD', 'foo');
const oD = std.param.Object('myObjectD', {
  s: 'bar',
  b: true,
  o: {
    xxx: 'yyy',
  },
});

std.log({
  myBoolean: bD,
  myNumber: nD,
  myString: sD,
  myObject: oD,
});
