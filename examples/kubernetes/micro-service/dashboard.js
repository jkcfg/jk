import * as G from '@jkcfg/grafana';

const r = '2m'; // Time window for range vectors.
const selector = service => `job='${service.name}'`;
const ServiceRPS = selector => `sum by (code)(sum(irate(http_requests_total{${selector}}[${r}])))`;
const ServiceLatency = selector => [
  `histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{${selector}}[${r}])) by (route) * 1e3`,
  `histogram_quantile(0.50, sum(rate(http_request_duration_seconds_bucket{${selector}}[${r}])) by (route) * 1e3`,
  `sum(rate(http_request_total{${selector}}[${r}])) / sum(rate(http_request_duration_seconds_count{${selector}}[${r}])) * 1e3`,
];


const RPSHttp = service => new G.Dashboard(`Service > ${service.name}`)
  .addPanel(
    new G.Graph(`${service.name} RPS`, {
      dataSource: '$PROMETHEUS_DS',
    })
      .addTargets([
        new G.Prometheus(ServiceRPS(selector(service)), { legendFormat: '{{code}}' }),
      ]),
    {
      gridPos: {
        x: 0, y: 0, w: 12, h: 7,
      },
    },
  )
  .addPanel(
    new G.Graph(`${service.name} Latency`, {
      dataSource: '$PROMETHEUS_DS',
      yAxis: [
        new G.YAxis({ format: 'ms' }),
        new G.YAxis(),
      ],
    })
      .addTargets([
        new G.Prometheus(ServiceLatency(selector(service))[0], { legendFormat: '99th percentile' }),
        new G.Prometheus(ServiceLatency(selector(service))[1], { legendFormat: 'median' }),
        new G.Prometheus(ServiceLatency(selector(service))[2], { legendFormat: 'mean' }),
      ]),
    {
      gridPos: {
        x: 12, y: 0, w: 12, h: 7,
      },
    },
  );

// In real life we probably want something different. I'd like to enable dynamic
// imports so we can import dashboards js definitions from string descriptions:
//   https://developers.google.com/web/updates/2017/11/dynamic-import
const dashboards = {
  'service.RPS.HTTP': RPSHttp,
};

function Dashboard(service) {
  if (!service.dashboards) {
    return [];
  }
  return service.dashboards.map(d => dashboards[d](service));
}

export {
  RPSHttp,
  Dashboard,
};
