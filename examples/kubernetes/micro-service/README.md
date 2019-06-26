# `micro-service.js`

An example generating various artefacts from a high level micro-service
definition. This example generates:

- `Namespace`, `Deployment`, `Service` and `Ingress` Kubernetes objects,
- A grafana dashboard stored in a `ConfigMap`,
- A Prometheus alert stored in a `PrometheusRule` custom resource, ready to be
  picked up by the [Prometheus operator][prom-operator].

There are two way to run this example:

1. The micro-service is defined with code in `billing.js`:

```console
$ npm install @jkcfg/kubernetes@0.2.1
$ npm install @jkcfg/grafana@0.1.0
$ jk generate -v billing.js
```

2. The micro-service is defined in a YAML file in `billing.yaml`:

```console
$ npm install @jkcfg/kubernetes@0.2.1
$ npm install @jkcfg/grafana@0.1.0
$ jk generate -v -f billing.yaml index.js
```

[prom-operator]: https://github.com/coreos/prometheus-operator
