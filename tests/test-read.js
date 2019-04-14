import std from '@jkcfg/std';

function writeErr(err) {
  std.write(`[ERROR] ${err.toString()}`);
}

// Read bytes from a file and write them exactly as found.
const raw = std.read('test-read.expected/raw.txt', { encoding: std.Encoding.Bytes });
raw.then(s => std.write(String.fromCharCode(...s), 'raw.txt', { format: std.Format.Raw }),
         writeErr);

// Read a UTF8 file into a JS string, and write it back.
const utf16 = std.read('test-read.expected/utf16.txt', { encoding: std.Encoding.String });
utf16.then(s => std.write(s, 'utf16.txt', { format: std.Format.Raw }),
           writeErr);

// Read a JSON file, implicitly as an object, modify it, then write it
// as a YAML file.
const json = std.read('foo.json');
json.then((s) => {
  if (typeof s !== 'object') {
    std.write(`[ERROR] value of read({ encoding: JSON }) is ${typeof s} instead of expected 'object'`);
    return;
  }
  const v = s;
  v.baz = 7;
  std.write(v, 'foo.json.json');
  std.write({ config: v }, 'foo.json.yaml');
}, writeErr);

// Read a YAML file as an object, and write modified objects back as JSON and YAML.
const yaml = std.read('foo.yaml');
yaml.then((s) => {
  if (typeof s !== 'object') {
    std.write(`[ERROR] value of read({ encoding: YAML }) is ${typeof s} instead of expected 'object'`);
    return;
  }
  std.write(s, 'foo.yaml.json');
  const v = s;
  v.config.baz = 7;
  std.write(v, 'foo.yaml.yaml');
}, writeErr);
