import std from '@jkcfg/std';

// test-log.js already has a bunch of tests, we just test write-specific things
// here, mainly writing to files.

const o = { kind: 'Bar', foo: { number: 1.2, string: 'mystring' } };
std.write(o, 'test-write.json');
std.write(o, 'test-write.yaml');
std.write(o, 'test-write.yml');
std.write(o, 'test-write-format-overwrite.yaml', { format: std.Format.JSON });

// Test the overwrite option: we don't write a file if it already exists and
// overwrite is false. overwrite defaults to true.
std.write(o, 'test-write-overwrite.json');
std.write(o, 'test-write-overwrite-no-file.json', { overwrite: false });
std.write({ ...o, foo: { ...o.foo, number: 1.3 } }, 'test-write-overwrite.json', { overwrite: false });
std.write({ ...o, foo: { ...o.foo, string: 'yourstring' } }, 'test-write-overwrite.json', { overwrite: true });

// Test writing a string in a JSON file does print the string as a JSON document.
std.write('success', 'test-write-json-string.json');
