import { RPCSync } from './internal/rpc';
import { valueFromUTF8Bytes } from './internal/data';

export function validate(obj: any, schema: Record<string, any>): 'ok' | string[] {
  const bytes = RPCSync('std.validate.schema', JSON.stringify(obj), JSON.stringify(schema));
  const results = valueFromUTF8Bytes(bytes);
  if (results === null) {
    return 'ok';
  }
  if (Array.isArray(results)) {
    return results;
  }
  throw new Error(`unexpected return value from RPC: ${results}`);
}
