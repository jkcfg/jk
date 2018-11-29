import std from 'std';

const r = std.read('test-read.js.expected');
r.then(std.write, err => std.write(`[ERROR] ${err.toString()}`));
