import { parse, print, Format } from '@jkcfg/std';

const json = parse('{ "json" : "ok" }', Format.JSON);
print(json);

const yaml = parse('yaml: ok', Format.YAML);
print(yaml);

const yamls = parse(`
---
foo: 1
---
bar: 2`, Format.YAMLStream);

print(yamls);
