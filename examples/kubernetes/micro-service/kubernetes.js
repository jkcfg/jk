import * as api from '@jkcfg/kubernetes/api';

function Namespace(service) {
  return new api.core.v1.Namespace(service.namespace);
}

function Deployment(service) {
  return new api.apps.v1.Deployment(service.name, {
    metadata: {
      namespace: service.namespace,
      labels: {
        app: service.name,
      },
    },
    spec: {
      replicas: service.replicas,
      revisionHistoryLimit: 2,
      strategy: {
        rollingUpdate: {
          maxUnavailable: 0,
          maxSurge: 1,
        },
      },
      template: {
        metadata: {
          labels: {
            app: service.name,
          },
        },
        spec: {
          containers: [{
            name: service.name,
            image: service.image,
            ports: [{
              containerPort: service.port,
            }],
          }],
        },
      },
    },
  });
}

function Service(service) {
  return new api.core.v1.Service(service.name, {
    metadata: {
      namespace: service.namespace,
      labels: {
        app: service.name,
      },
    },
    spec: {
      selector: {
        app: service.name,
      },
      ports: [{
        port: service.port,
      }],
    },
  });
}

function Ingress(service) {
  return new api.extensions.v1beta1.Ingress(service.name, {
    metadata: {
      namespace: service.namespace,
      annotations: {
        'nginx.ingress.kubernetes.io/rewrite-target': '/',
      },
    },
    spec: {
      rules: [{
        http: {
          paths: [{
            path: service.ingress.path,
            backend: {
              serviceName: service.name,
              servicePort: service.port,
            },
          }],
        },
      }],
    },
  });
}

export {
  Deployment,
  Ingress,
  Namespace,
  Service,
};
