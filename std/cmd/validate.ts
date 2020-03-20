import { log, Format } from '../index';
import * as host from '@jkcfg/std/internal/host'; // magic module
import * as param from '../param';
import { formatError, normaliseResult, ValidationError, ValidationResult, ValidateFnResult } from '../validation';
import { valuesFormatFromPath } from '../read';

export type ValidateFn = (obj: any) => ValidateFnResult | Promise<ValidateFnResult>;

interface FileResult {
  path: string;
  result: ValidationResult;
}

function reduce(results: ValidationResult[]): ValidationResult {
  return results.reduce((a: ValidationResult, b: ValidationResult): ValidationResult => {
    if (a == 'ok') return b;
    if (b == 'ok') return a;
    return Array.prototype.concat(a, b);
  }, 'ok');
}

export default function validate(fn: ValidateFn): void {
  const inputFiles = param.Object('jk.validate.input', {});
  const files = Object.keys(inputFiles);

  function validateValue(v: any): Promise<ValidationResult> {
    return Promise.resolve(fn(v)).then(normaliseResult);
  }

  async function validateFile(path: string): Promise<FileResult> {
    const format = valuesFormatFromPath(path);
    const obj = await host.read(path, { format });
    switch (format) {
    case Format.YAMLStream:
    case Format.JSONStream:
      const results: Promise<ValidationResult>[] = obj.map(validateValue);
      const resolvedResults = await Promise.all(results);
      return { path, result: reduce(resolvedResults) };
    default:
      const result = await validateValue(obj);
      return { path, result };
    }
  }

  const objects = files.map(validateFile);
  Promise.all(objects).then((results): void => {
    for (const { path, result } of results) {
      if (result === 'ok') {
        log(`${path}: ok`);
      } else {
        for (const err of result) {
          log(formatError(path, err));
        }
      }
    }
  });
}
