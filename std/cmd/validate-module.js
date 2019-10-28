import validate from '@jkcfg/std/cmd/validate';
import fn from '%s';

if (typeof fn !== 'function') {
  throw new Error('default export of given module is not a function');
}

validate(fn);
