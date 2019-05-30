import { parse, Format, log } from '@jkcfg/std';

for (const { string, format, message } of [
  { string: '"anything"', format: Format.FromExtension, message: 'FromExtension format' },
  { string: 'garbage', format: Format.JSON, message: 'Garbage JSON' },
  { string: '1\n2\n', format: Format.JSON, message: 'Multiple JSON values' },
  { string: '{{ invalid }}', format: Format.YAML, message: 'Invalid YAML' },
  { string: '1\nfoo\n', format: Format.JSONStream, message: 'Invalid JSON stream' },
]) {
  try {
    const obj = parse(string, format);
    log(obj);
  } catch (_) {
    log(`${message} correctly errored.`);
  }
}
