import { parse, log, Format } from '@jkcfg/std';

const json = parse('{ "json" : "ok" }', Format.JSON);
log(json);

const yaml = parse('yaml: ok', Format.YAML);
log(yaml);

const yamls = parse(`
---
foo: 1
---
bar: 2`, Format.YAMLStream);

log(yamls);
