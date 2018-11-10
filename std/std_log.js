import { write } from 'std_write';

function log(value, format) {
  write(value, '', format);
}

export {
  log,
};
