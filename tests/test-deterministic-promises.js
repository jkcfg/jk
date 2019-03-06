import std from '@jkcfg/std';

const resolvedOrder = [];
const promises = [];
for (let i = 0; i < 100; i += 1) {
  promises.push(std.read('success.json').then(() => resolvedOrder.push(i)));
}

Promise.all(promises).then(() => std.log(resolvedOrder.join(' ')));
