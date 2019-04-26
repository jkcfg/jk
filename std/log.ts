import { write, WriteOptions } from './write';

function log(value: any, options?: WriteOptions): void {
  write(value, '', options);
}

export {
  log,
};
