// merge provides a small set of ways to transform values (usually
// objects).

// patch returns a new value that has the fields of `obj`, except
// where overridden by fields in `patchObj`. Entries in common are
// themselves patched.
function patch(obj, patchObj) {
  switch (typeof obj) {
  case 'object': {
    const result = {};
    for (const [k, v] of Object.entries(obj)) {
      if (k in patchObj) {
        result[k] = patch(v, patchObj[k]);
      } else {
        result[k] = v;
      }
    }
    for (const [pk, pv] of Object.entries(patchObj)) {
      if (!(pk in obj)) {
        result[pk] = pv;
      }
    }
    return result;
  }
  case 'string':
  case 'number':
  case 'boolean':
    return patchObj;
  default:
    throw new Error(`unhandled patch case: ${typeof obj}`);
  }
}

// Interpret a series of transformations expressed either as object
// patches (as in the argument to `patch` in this module), or
// functions. Usually the first argument will be an object,
// representing an initial value, but it can be a function (that will
// be given an empty object as its argument).
function mix(...transforms) {
  let r = {};

  for (const transform of transforms) {
    switch (typeof transform) {
    case 'object':
      r = patch(r, transform);
      break;
    case 'function':
      r = transform(r);
      break;
    default:
      throw new TypeError('only objects and functions allowed as arguments');
    }
  }

  return r;
}

export { mix, patch };
