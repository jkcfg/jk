import * as param from '@jkcfg/std/param';
import { generate, OutputFormat } from '@jkcfg/std/cmd/generate';
import generateDefinition from '%s';

let inputParams = {
  stdout: param.Boolean('jk.generate.stdout', false),
};

let format = param.String('jk.generate.format', undefined);
switch (format) {
case "json":
  inputParams.format = OutputFormat.JSON;
  break;
case "yaml":
  inputParams.format = OutputFormat.YAML;
  break;
default:
  break;
}

generate(generateDefinition, inputParams);
