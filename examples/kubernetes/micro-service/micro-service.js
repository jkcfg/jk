import * as param from '@jkcfg/std/param';
import {
  Namespace, Deployment, Service, Ingress,
} from './kubernetes';

const service = param.Object('service');
const ns = service.namespace;

export default [
  { file: `${ns}-ns.yaml`, value: Namespace(service) },
  { file: `${ns}/${service.name}-deploy.yaml`, value: Deployment(service) },
  { file: `${ns}/${service.name}-svc.yaml`, value: Service(service) },
  { file: `${ns}/${service.name}-ingress.yaml`, value: Ingress(service) },
];
