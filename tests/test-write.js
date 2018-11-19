import std from 'std';

// test-log.js already has a bunch of tests, we just test write-specific things
// here, mainly writing to files.

const o = { kind: 'Bar', foo: { number: 1.2, string: 'mystring' } };
std.write(o, 'test-write.json');
std.write(o, 'test-write.yaml');
std.write(o, 'test-write.yml');
std.write(o, 'test-write-override.yaml', { format: std.Format.JSON });
