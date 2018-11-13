import { write } from 'std_write';

function log(value, options) {
  write(value, '', options);
}

export {
  log,
};
