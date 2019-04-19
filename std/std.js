import { log } from 'std_log';
import { Format, write } from 'std_write';
import { Encoding, read } from 'std_read';
import { info, dir } from 'std_fs';
import { param } from 'std_param';
import { mix, patch, merge } from 'std_merge';

// The default export is deprecated and will be removed in 3.0.0.

export default {
  log,
  Format,
  write,
  Encoding,
  read,
  info,
  dir,
  param,
  mix,
  patch,
  merge,
};

export { log } from 'std_log';
export { Format, write } from 'std_write';
export { Encoding, read } from 'std_read';
export { info, dir } from 'std_fs';
export { param } from 'std_param';
export { mix, patch, merge } from 'std_merge';
