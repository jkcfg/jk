import * as param from '@jkcfg/std/param';
import { generate, OutputFormat, maybeSetFormat } from '@jkcfg/std/cmd/generate';
import generateDefinition from '%s';

let inputParams = {
  stdout: param.Boolean('jk.generate.stdout', false),
};
maybeSetFormat(inputParams, param.String('jk.generate.format', undefined));

generate(generateDefinition, inputParams);
