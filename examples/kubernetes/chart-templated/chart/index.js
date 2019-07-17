import { generateChart, loadModuleTemplates } from '@jkcfg/kubernetes/chart';
import handlebars from 'handlebars/lib/handlebars';
import * as resource from '@jkcfg/std/resource';

const templates = loadModuleTemplates(handlebars.compile, resource);
const defaults = resource.read('./defaults.yaml');

export default (paramMod => generateChart(templates, defaults, paramMod));
