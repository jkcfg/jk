import { stringify, Format } from '@jkcfg/std';
import * as k from './kubernetes';
import { Dashboard } from './dashboard';
import { PrometheusRule } from './alert';

export function MicroService(service) {
  const ns = service.namespace;

  const dashboards = {
    dashboard: stringify(Dashboard(service), Format.JSON),
  };

  return [
    { path: `${ns}-ns.yaml`, value: k.Namespace(service) },
    { path: `${ns}/${service.name}-deploy.yaml`, value: k.Deployment(service) },
    { path: `${ns}/${service.name}-svc.yaml`, value: k.Service(service) },
    { path: `${ns}/${service.name}-ingress.yaml`, value: k.Ingress(service) },
    { path: `${ns}/${service.name}-dashboards-cm.yaml`, value: k.ConfigMap(service, `${service.name}-dashboards`, dashboards) },
    { path: `${ns}/${service.name}-prometheus-rule.yaml`, value: PrometheusRule(service) },
  ];
}
