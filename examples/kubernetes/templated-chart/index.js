// Example of a Helm chart analogue, using handlebars

import { generateChart, loadModuleTemplates } from '@jkcfg/kubernetes/chart';
import * as resource from '@jkcfg/std/resource';
import * as param from '@jkcfg/std/param';
import handlebars from 'handlebars/lib/handlebars';

const defaults = {
  name: 'helloworld',
  app: 'hello',
  image: {
    repository: 'weaveworks/helloworld',
    tag: 'v1'
  }
};

const templates = loadModuleTemplates(handlebars.compile, resource);

export default generateChart(templates, defaults, param);
