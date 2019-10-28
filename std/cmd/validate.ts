import * as std from '../index';
import * as param from '../param';

// ValidateResult is the canonical type of results for a validation
// procedure.
type ValidateResult = 'ok' | string[];

// ValidateFnResult is the range of results we accept from an ad-hoc
// procedure given to us.
type ValidateFnResult = boolean | string | string[];
type ValidateFn = (obj: any) => ValidateFnResult;

function normaliseResult(result: ValidateFnResult): ValidateResult {
  switch (typeof result) {
  case 'string':
    if (result === 'ok') return result;
    return [result];
  case 'boolean':
    if (result) return 'ok';
    return ['value not valid'];
  case 'object':
    if (Array.isArray(result)) return result;
    break;
  default:
  }
  throw new Error(`unrecognised result from validation function: ${result}`);
}

interface FileResult {
  path: string;
  result: ValidateResult;
}

export default function validate(fn: ValidateFn): void {
  const inputFiles = param.Object('jk.validate.input', {});
  const files = Object.keys(inputFiles);

  const validateFile = async function vf(path: string): Promise<FileResult> {
    const obj = await std.read(path);
    return { path, result: normaliseResult(fn(obj)) };
  };

  const objects = files.map(validateFile);
  Promise.all(objects).then((results): void => {
    for (const { path, result } of results) {
      if (result === 'ok') {
        std.log(`${path}: ok`);
      } else {
        std.log(`${path}:`);
        for (const err of result) {
          std.log(`  error: ${err}`);
        }
      }
    }
  });
}
