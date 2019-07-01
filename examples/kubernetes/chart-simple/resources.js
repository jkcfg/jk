import { core, apps } from '@jkcfg/kubernetes/api';

function resources(Values) {
  return [
    new apps.v1.Deployment(`${Values.name}-dep`, {
      spec: {
        template: {
          labels: { app: Values.app },
          spec: {
            containers: {
              hello: {
                image: `${Values.image.repository}:${Values.image.tag}`,
              },
            },
          },
        },
      },
    }),
    new core.v1.Service(`${Values.name}-svc`, {
      metadata: {
        labels: { app: Values.app },
      },
      spec: {
        selector: {
          app: Values.app,
        },
      },
    })];
}

export default resources;
