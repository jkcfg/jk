export default function doTest(output) {
  // We can print basic types.
  output(1.2);
  output('foo');
  output(undefined);
  output(null);

  // We can print an object to stderr.
  output({ kind: 'Bar', foo: 1.2 });
}
