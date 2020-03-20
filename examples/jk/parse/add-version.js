import { parse, stringify, Format } from '@jkcfg/std';
import { merge } from '@jkcfg/std/merge';

// We're going to extract embedded JSON from a YAML resource, alter
// it, then reconstruct the resource.

export default function addversion(resource) {
  const { data, ...rest } = resource;
  const inConfig = parse(data['config.json'], Format.JSON);
  const outConfig = merge(inConfig, { version: 2 });
  data['config.json'] = stringify(outConfig, Format.JSON);
  return { data, ...rest };
}
