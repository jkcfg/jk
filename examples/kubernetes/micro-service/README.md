# `micro-service.js`

An example generating `Namespace`, `Deployment`, `Service` and `Ingress`
Kubernetes objects from a high level micro-service definition stored as input
parameters in `billing.yaml`.

Run this example with:

```console
$ npm install @jkcfg/kubernetes
$ jk generate -f billing.yaml micro-service.js
```
