import { unparse, Format, log } from '@jkcfg/std';

log('# JSON');
const json = unparse({ json: 'ok' }, Format.JSON);
log(json);

log('# YAML');
const yaml = unparse({ yaml: 'ok' }, Format.YAML);
log(yaml);

log('# JSON stream');
const jsons = unparse([{ json: 1 }, { json: 2 }], Format.JSONStream);
log(jsons);

log('# YAML stream');
const yamls = unparse([{ yaml: 1 }, { yaml: 2 }], Format.YAMLStream);
log(yamls);
