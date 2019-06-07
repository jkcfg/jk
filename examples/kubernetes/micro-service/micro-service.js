import { unparse, Format } from '@jkcfg/std';
import * as param from '@jkcfg/std/param';
import {
  Namespace, Deployment, Service, Ingress, ConfigMap,
} from './kubernetes';
import { Dashboard } from './dashboard';

const service = param.Object('service');
const ns = service.namespace;

const dashboards = {
  dashboard: unparse(Dashboard(service), Format.JSON),
};

export default [
  { file: `${ns}-ns.yaml`, value: Namespace(service) },
  { file: `${ns}/${service.name}-deploy.yaml`, value: Deployment(service) },
  { file: `${ns}/${service.name}-svc.yaml`, value: Service(service) },
  { file: `${ns}/${service.name}-ingress.yaml`, value: Ingress(service) },
  { file: `${ns}/${service.name}-dashboards-cm.yaml`, value: ConfigMap(service, `${service.name}-dashboards`, dashboards) },
];
