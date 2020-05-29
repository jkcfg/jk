import { Format } from '@jkcfg/std';

const array = [
  { message: 'hello' },
];

export default [
  { format: Format.JSON, path: 'jsonarray.json', value: array },
  { format: Format.JSONStream, path: 'jsonstream.json', value: array },
  { format: Format.YAML, path: 'yamlarray.yaml', value: array },
  { format: Format.YAMLStream, path: 'yamlstream.yaml', value: array },
];
