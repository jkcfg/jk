import * as prometheus from './prometheus';

const r = '2m';
const selector = service => `job=${service.name}`;
const ErrorRate = selector => `rate(http_request_total{${selector},code=~"5.."}[${r}])
    / rate(http_request_duration_seconds_count{${selector}}[${r}])`;

function RPSHTTPHighErrorRate(service) {
  return {
    alert: 'HighErrorRate',
    expr: `${ErrorRate(selector(service))} * 100 > 10`,
    for: '5m',
    labels: {
      severity: 'critical',
    },
    annotations: {
      service: service.name,
      description: `More than 10% of requests to the ${service.name} service are failing with 5xx errors`,
      details: '{{$value | printf "%.1f"}}% errors for more than 5m',
    },
  };
}

// In real life we probably want something different. I'd like to enable dynamic
// imports so we can import dashboards js definitions from string descriptions:
//   https://developers.google.com/web/updates/2017/11/dynamic-import
const alerts = {
  'service.RPS.HTTP.HighErrorRate': RPSHTTPHighErrorRate,
};

function rules(service) {
  if (!service.alerts) {
    return [];
  }

  return service.alerts.map(a => alerts[a](service));
}

function PrometheusRule(service) {
  return new prometheus.PrometheusRule(service.name, {
    metadata: {
      labels: {
        app: service.name,
        maintainer: service.maintainer,
        prometheus: 'global',
        role: 'alert-rules',
      },
    },
    spec: {
      groups: [{
        name: `${service.name}-alerts.rules`,
        rules: rules(service),
      }],
    },
  });
}

export {
  PrometheusRule,
};
