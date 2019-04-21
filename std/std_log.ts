import { write, WriteOptions } from './std_write';

function log(value: any, options?: WriteOptions): void {
  write(value, '', options);
}

export {
  log,
};
