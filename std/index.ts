import { log } from './internal/log';
import { Format, write } from './internal/write';
import { Encoding, read } from './internal/read';
import { info, dir } from './fs';
import { param } from './param';
import { mix, patch, merge } from './merge';

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

export { log } from './internal/log';
export { Format, write } from './internal/write';
export { Encoding, read } from './internal/read';
export { info, dir } from './fs';
export { param } from './param';
export { mix, patch, merge } from './merge';
export { parse, unparse } from './parse';
