import { overlay } from '@jkcfg/kubernetes/overlay';
import { valuesForGenerate } from '@jkcfg/kubernetes/generate';
import * as param from '@jkcfg/std/param';

// This is similar to the first example, but uses an object rather
// than going straight to the filesystem, and overlays further
// changes.
//
// The `bases` part loads and interprets kustomization files, so
// another way to do the `overlay-simple` example would be:
//
//     overlay('.', { bases: ['.'] });
//
const kustom = {
  bases: ['.'],
  resources: ['service.yaml'],
  commonLabels: {
    team: 'strange',
  },
  patches: ['service-selector.json'],
  namespace: param.String('namespace', 'default'),
};

export default valuesForGenerate(overlay('.', kustom));
