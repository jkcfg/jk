# `micro-service.js`

An example generating various artefacts from a high level micro-service
definition stored as input parameters in `billing.yaml`:

- `Namespace`, `Deployment`, `Service` and `Ingress` Kubernetes objects,
- A grafana dashboard stored in a `ConfigMap`,
- A Prometheus alert stored in a `PrometheusRule` custom resource, ready to be
  picked up by the [Prometheus operator][prom-operator].

Run this example with:

```console
$ npm install @jkcfg/kubernetes
$ jk generate -f billing.yaml micro-service.js
```

[prom-operator]: https://github.com/coreos/prometheus-operator
