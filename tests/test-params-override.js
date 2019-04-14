import std from '@jkcfg/std';

const b = std.param.Boolean('myBoolean', false);
const n = std.param.Number('myNumber', 3.14);
const s = std.param.String('myString', 'foo');
const o = std.param.Object('myObject', {
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
