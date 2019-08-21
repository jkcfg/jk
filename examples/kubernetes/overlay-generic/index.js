import * as param from '@jkcfg/std/param';
import { merge, deepWithKey } from '@jkcfg/std/merge';
import { overlay } from '@jkcfg/kubernetes/overlay';
import { core } from '@jkcfg/kubernetes/api';
import { valuesForGenerate } from '@jkcfg/kubernetes/generate';

import { resources as submoduleResources } from './submodule';

// This transformation adds a sidecar to each Deployment
function addSidecar(maybeDeployment) {
  if (maybeDeployment.kind !== 'Deployment') {
    return maybeDeployment;
  }

  return merge(maybeDeployment, {
    spec: {
      template: {
        spec: {
          containers: [{ name: 'sidecar', image: 'foobar:1.1.0' }],
        },
      },
    },
  }, {
    spec: {
      template: {
        spec: {
          containers: deepWithKey('name'),
        },
      },
    },
  });
}

// The generate function is parameterised by the name of an
// environment -- it's just a string -- which we'll use to specialise
// the generated resources.
function generateEnv(env) {
  const ns = `${env}-env`;
  const nsResource = new core.v1.Namespace(ns, {});

  return overlay('.', {
    namespace: ns,
    commonLabels: { env },
    generatedResources: [[nsResource], submoduleResources()],
    transformations: [addSidecar],
  });
}

export default generateEnv(param.String('env', 'dev')).then(valuesForGenerate);
