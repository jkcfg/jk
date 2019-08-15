import { log } from './log';
import { Format, write } from './write';
import { Encoding, read } from './read';
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

export { log } from './log';
export { Format, write } from './write';
export { Encoding, read } from './read';
export { parse, unparse } from './parse';
