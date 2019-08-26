import { valueFromUTF8Bytes } from './data';
import { RPC } from './rpc';

export function echo(...args) {
  return RPC("debug.echo", ...args).then(valueFromUTF8Bytes);
}
