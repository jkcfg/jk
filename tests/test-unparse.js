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

log('# HCL config');
const config = {
  provider: {
    github: {
      organization: 'myorg',
    },
  },
  github_membership: {
    myorg_foo: {
      username: 'foo',
      role: 'admin',
    },
  },
};
const hcl = unparse(config, Format.HCL);
log(hcl);

log('# Unsupported format');
try {
  const str = unparse({ foo: 2 }, Format.FromExtension);
  log(str);
} catch (_) {
  log('Unsupported format correctly errored.');
}
