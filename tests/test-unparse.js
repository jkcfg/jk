import { stringify, Format, log } from '@jkcfg/std';

log('# JSON');
const json = stringify({ json: 'ok' }, Format.JSON);
log(json);

log('# YAML');
const yaml = stringify({ yaml: 'ok' }, Format.YAML);
log(yaml);

log('# JSON stream');
const jsons = stringify([{ json: 1 }, { json: 2 }], Format.JSONStream);
log(jsons);

log('# YAML stream');
const yamls = stringify([{ yaml: 1 }, { yaml: 2 }], Format.YAMLStream);
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
const hcl = stringify(config, Format.HCL);
log(hcl);

log('# Unsupported format');
try {
  const str = stringify({ foo: 2 }, Format.FromExtension);
  log(str);
} catch (_) {
  log('Unsupported format correctly errored.');
}
