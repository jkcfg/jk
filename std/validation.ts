/**
 * @module std/validation
 */

interface Location {
  line: number;
  column: number;
}

/**
 * ValidationError represents a specific problem encountered when
 * validating a value.
 */
export interface ValidationError {
  msg: string;
  path?: string;
  start?: Location;
  end?: Location;
}

function isValidationError(err: string | ValidationError): err is ValidationError {
  return (typeof err === 'object' && 'msg' in err);
}

// ValidationResult is the canonical type of results for a validation
// procedure.
export type ValidationResult = 'ok' | ValidationError[];

// ValidateFnResult is the range of results we accept from an ad-hoc
// procedure given to us.
export type ValidateFnResult = boolean | string | (string | ValidationError)[];

function normaliseError(err: string | ValidationError): ValidationError {
  if (typeof err === 'string') return { msg: err };
  if (typeof err === 'object' && isValidationError(err)) return err;
  throw new Error(`unrecognised result from validation function: ${err}`);
}

// normaliseResult maps from the accepted return values of a
// validation function to the more restrictive `ValidateResult`, so
// that results can be dealt with uniformly.
export function normaliseResult(result: ValidateFnResult): ValidationResult {
  switch (typeof result) {
  case 'string':
    if (result === 'ok') return result;
    return [{ msg: result }];
  case 'boolean':
    if (result) return 'ok';
    return [{ msg: 'value not valid' }];
  case 'object':
    if (Array.isArray(result)) return result.map(normaliseError);
    if (isValidationError(result)) return [result];
    break;
  default:
  }
  throw new Error(`unrecognised result from validation function: ${result}`);
}

/**
 * formatError formats a validation error in a standard way. Since the
 * errors do not include the source, this must be supplied.
 */
export function formatError(source: string, err: ValidationError): string {
  let out = err.msg;
  if (err.path !== undefined) {
    out = `${err.msg} at ${err.path}`;
  }
  if (err.start !== undefined) {
    if (err.end !== undefined) {
      out = `${err.start.line}.${err.start.column}-${err.end.line}.${err.end.column}: ${out}`;
    } else {
      out = `${err.start.line}.${err.start.column}: ${out}`;
    }
  } else { // if no location, put a spacer between the source file and message
    out = ` ${out}`;
  }
  return `${source}:${out}`;
}
