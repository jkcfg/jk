import std from 'std';

const r = std.read('test-read.js.expected');
r.then(s => std.write(s, '', { format: std.Format.Raw }), err => std.write(`[ERROR] ${err.toString()}`));
