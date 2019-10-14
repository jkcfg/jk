import { RPC, RPCSync } from './internal/rpc';
import { valueFromUTF8Bytes } from './internal/data';

type Result = 'ok' | string[];

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

export function validateBySchema(obj: any, schema: Record<string, any>): Result {
  return decodeResponse(RPCSync('std.validate.schema', JSON.stringify(obj), JSON.stringify(schema)));
}

export function validateByFile(obj: any, path: string): Promise<Result> {
  return RPC('std.validate.schemafile', JSON.stringify(obj), path, '').then(decodeResponse);
}

// This is intended to be used by invoking `@jkcfg/std/resource#withModuleRef`
export function validateByResource(obj: any, path: string, moduleRef: string): Promise<Result> {
  return RPC('std.validate.schemafile', JSON.stringify(obj), path, moduleRef).then(decodeResponse);
}
