import std from '@jkcfg/std';

function copy(...filenames) {
  for (const filename of filenames) {
    std.read(`test-issue-0071/${filename}`, { encoding: std.Encoding.String }).then(
      content => (std.write(content, filename, { format: std.Format.Raw })),
      err => std.write(`[ERROR] ${err.toString()}`),
    );
  }
}

copy(
  '.editorconfig',
  'LICENSE',
);
