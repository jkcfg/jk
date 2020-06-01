import { Format } from '@jkcfg/std';

function valueAndFormat(f) {
  return {
    format: f,
    value: [
      { item1: Format[f] },
      { item2: Format[f] },
      { item3: Format[f] },
    ],
  };
}

export default [
  { path: 'jsonarray.json', ...valueAndFormat(Format.JSON) },
  { path: 'jsonstream.json', ...valueAndFormat(Format.JSONStream) },
  { path: 'yamlarray.yaml', ...valueAndFormat(Format.YAML) },
  { path: 'yamlstream.yaml', ...valueAndFormat(Format.YAMLStream) },
];
