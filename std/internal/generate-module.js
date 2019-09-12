import * as param from '@jkcfg/std/param';
import { generate } from '@jkcfg/std/generate';
import generateDefinition from '%s';

const inputParams = {
  stdout: param.Boolean('jk.generate.stdout', false),
};

generate(generateDefinition, inputParams);
