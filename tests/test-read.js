import std from 'std';

const r = std.read('here is the URL');
r.then(str => std.write(`[OUT] ${str}`), err => std.write(`[ERROR] ${err.toString()}`));
