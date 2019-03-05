import std from '@jkcfg/std';

const resolvedOrder = [];
const checkRead = i => std.read('success.json', { encoding: std.Encoding.String }).then(() => resolvedOrder.push(i));

const promises = [];
for (let i = 0; i < 100; i += 1) {
  promises.push(checkRead(i));
}

Promise.all(promises).then(() => std.log(resolvedOrder.join(' ')));
