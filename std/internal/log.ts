import { write, WriteOptions } from './write';

function log(value: any, options?: WriteOptions): void {
  if (value === undefined) {
    V8Worker2.print('undefined');
    return;
  }
  write(value, '', options);
}

export {
  log,
};
