/**
 * @module std
 */

import { write, WriteOptions } from './write';

export function log(value: any, options?: WriteOptions): void {
  if (value === undefined) {
    V8Worker2.print('undefined');
    return;
  }
  write(value, '', options);
}
