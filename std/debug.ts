/**
 * @module std/debug
 *
 * debug exists to help with testing the runtime.
 */

import { valueFromUTF8Bytes } from './internal/data';
import { RPC, RPCSync } from './internal/rpc';

export function echo(...args: any[]): Promise<any[]> {
  return RPC('debug.echo', ...args).then(valueFromUTF8Bytes);
}

export function echoSync(...args: any[]): any[] {
  return valueFromUTF8Bytes(RPCSync('debug.echo', ...args));
}
