import { RPC } from './rpc';

export function debug(...args) {
  return RPC("debug", ...args);
}
