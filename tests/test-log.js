import std from 'std';

// We can print an object to stdout.
std.log({kind: "Bar", foo: 1.2});

// Key order is deterministic.
std.log({foo: 1.2, kind: "Bar"});
