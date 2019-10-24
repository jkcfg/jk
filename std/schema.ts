/**
 * @module std/schema
 */

import { RPC, RPCSync } from './internal/rpc';
import { valueFromUTF8Bytes } from './internal/data';

// These two types could be put in a more general validation module
// (along with convenience formatters), since they are generic.

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

export type Result = 'ok' | ValidationError[];

function decodeResponse(bytes: Uint8Array): Result {
  const results = valueFromUTF8Bytes(bytes);
  if (results === null) {
    return 'ok';
  }
  if (Array.isArray(results)) {
    return results;
  }
  throw new Error(`unexpected return value from RPC: ${results}`);
}

/**
 * validateWithObject validates a value using a JSON schema supplied
 * as object.
 *
 * ```typescript
 * const result = validateWithObject(5, { type: 'number' });
 * ```
 */
export function validateWithObject(obj: any, schema: Record<string, any>): Result {
  return decodeResponse(RPCSync('std.validate.schema', JSON.stringify(obj), JSON.stringify(schema)));
}

/**
 * validateWithFile validates a value using a schema located at
 * the path (relative to the input directory).
 */
export function validateWithFile(obj: any, path: string): Promise<Result> {
  return RPC('std.validate.schemafile', JSON.stringify(obj), path, '').then(decodeResponse);
}

/**
 * validateWithResource validates a value using a schema location at
 * the given path relative to the module represented by `moduleRef`;
 * this is intended to be used by wrapping it in
 * `@jkcfg/std/resource#withModuleRef`.
 *
 * ```javascript
 * import { withModuleRef } from '@jkcfg/std/resource';
 *
 * export function validate(value) {
 *   return withModuleRef(ref => validateWithSchema(value, 'schema.json', ref));
 * }
 * ```
 */
export function validateWithResource(obj: any, path: string, moduleRef: string): Promise<Result> {
  return RPC('std.validate.schemafile', JSON.stringify(obj), path, moduleRef).then(decodeResponse);
}
