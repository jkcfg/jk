import std from 'std';

// test-log.js already has a bunch of tests, we just test write-specific things
// here, mainly writing to files.

const o = { kind: 'Bar', foo: { number: 1.2, string: 'mystring' } };
std.write(o, 'test-write.json');
std.write(o, 'test-write.yaml');
std.write(o, 'test-write.yml');
std.write(o, 'test-write-format-override.yaml', { format: std.Format.JSON });

// Test the override option: we don't write a file if it already exists and
// override is false. override defaults to true.
std.write(o, 'test-write-override.json');
std.write(o, 'test-write-override-no-file.json', { override: false });
std.write({ ...o, foo: { ...o.foo, number: 1.3 } }, 'test-write-override.json', { override: false });
std.write({ ...o, foo: { ...o.foo, string: 'yourstring' } }, 'test-write-override.json', { override: true });
