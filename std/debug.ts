import { valueFromUTF8Bytes } from './internal/data';
import { RPC } from './internal/rpc';

export function echo(...args: any[]): Promise<any[]> {
  return RPC('debug.echo', ...args).then(valueFromUTF8Bytes);
}
