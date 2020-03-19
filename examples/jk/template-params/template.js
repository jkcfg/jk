import * as param from '@jkcfg/std/param';
import { print, Format } from '@jkcfg/std';

const selector = param.Object('selector', false);

const obj = {
  apiVersion: 'v1',
  kind: 'Service',
  metadata: {
    name: param.String('name', 'service'),
    namespace: param.String('namespace', 'default'),
  },
  spec: {
    ports: [
      { port: param.Number('port', 8080) },
    ],
    selector: selector || { app: 'app' },
  },
};

print(obj, { format: Format.YAML });
