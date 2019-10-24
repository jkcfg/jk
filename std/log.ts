/**
 * @module std
 */

export function log(value: any): void {
  if (value === undefined) {
    V8Worker2.log('undefined');
    return;
  }
  if (typeof value === 'string') {
    V8Worker2.log(value);
    return;
  }
  V8Worker2.log(JSON.stringify(value));
}
