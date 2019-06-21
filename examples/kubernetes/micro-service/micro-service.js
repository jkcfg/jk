import { unparse, Format } from '@jkcfg/std';
import * as param from '@jkcfg/std/param';
import * as k from './kubernetes';
import { Dashboard } from './dashboard';
import { PrometheusRule } from './alert';

const service = param.Object('service');
const ns = service.namespace;

const dashboards = {
  dashboard: unparse(Dashboard(service), Format.JSON),
};

export default [
  { file: `${ns}-ns.yaml`, value: k.Namespace(service) },
  { file: `${ns}/${service.name}-deploy.yaml`, value: k.Deployment(service) },
  { file: `${ns}/${service.name}-svc.yaml`, value: k.Service(service) },
  { file: `${ns}/${service.name}-ingress.yaml`, value: k.Ingress(service) },
  { file: `${ns}/${service.name}-dashboards-cm.yaml`, value: k.ConfigMap(service, `${service.name}-dashboards`, dashboards) },
  { file: `${ns}/${service.name}-prometheus-rule.yaml`, value: PrometheusRule(service) },
];
