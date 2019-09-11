import transform from '@jkcfg/std/transform';
import fn from '%s';

if (typeof fn !== 'function') {
  throw new Error('default export of given module is not a function');
}

transform(fn);
