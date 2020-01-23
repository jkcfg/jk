import * as std from '../index';
import * as param from '../param';
import { formatError, normaliseResult, ValidationError, ValidationResult, ValidateFnResult } from '../validation';

export type ValidateFn = (obj: any) => ValidateFnResult | Promise<ValidateFnResult>;

interface FileResult {
  path: string;
  result: ValidationResult;
}

export default function validate(fn: ValidateFn): void {
  const inputFiles = param.Object('jk.validate.input', {});
  const files = Object.keys(inputFiles);

  const validateFile = async function vf(path: string): Promise<FileResult> {
    const obj = await std.read(path);
    const result = normaliseResult(await Promise.resolve(fn(obj)));
    return { path, result };
  };

  const objects = files.map(validateFile);
  Promise.all(objects).then((results): void => {
    for (const { path, result } of results) {
      if (result === 'ok') {
        std.log(`${path}: ok`);
      } else {
        for (const err of result) {
          std.log(formatError(path, err));
        }
      }
    }
  });
}
