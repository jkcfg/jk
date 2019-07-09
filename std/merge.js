/**
 * @module std/merge
 */

// patch returns a new value that has the fields of `obj`, except
// where overridden by fields in `patchObj`. Entries in common are
// themselves patched. This is similar to `merge` below, but always
// does a deep merge on objects, and always replaces other values.
export function patch(obj, patchObj) {
  switch (typeof obj) {
  case 'object': {
    const result = {};
    for (const [k, v] of Object.entries(obj)) {
      if (k in patchObj) {
        if (Array.isArray(v)) {
          result[k] = patchObj[k];
        } else {
          result[k] = patch(v, patchObj[k]);
        }
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

// merge returns a new value which is `a` merged additively with
// `b`. For values other than objects, this means addition (or
// concatenation), with a coercion if necessary.
//
// If both `a` and `b` are objects, there is some fine control over
// each field. If the key in `b` ends with a `+`, the values are
// summed; otherwise, the value is replaced. Any fields in
// `a` and not in `b` are also assigned in the result.
export function merge(a, b) {
  const [typeA, typeB] = [typeof a, typeof b];
  if (typeA === 'string') {
    if (typeB === 'string') {
      return a + b;
    }
    return a + JSON.stringify(b);
  }
  if (typeB === 'string') {
    return JSON.stringify(a) + b;
  }

  if (typeA === 'number' && typeB === 'number') return a + b;
  if (Array.isArray(a) && Array.isArray(b)) return [...a, ...b];
  if (typeA === 'object' && typeB === 'object') return objectMerge(a, b);

  throw new Error(`merge cannot combine values of types ${typeA} and ${typeB}`);
}

function objectMerge(obj, mergeObj) {
  const r = {};

  Object.assign(r, obj);
  for (let [key, value] of Object.entries(mergeObj)) {
    if (key.endsWith('+')) {
      key = key.slice(0, -1);
      if (key in obj) {
        r[key] = merge(obj[key], value);
        continue;
      }
    }
    r[key] = value;
  }
  return r;
}

// Interpret a series of transformations expressed either as object
// patches (as in the argument to `patch` in this module), or
// functions. Usually the first argument will be an object,
// representing an initial value, but it can be a function (that will
// be given an empty object as its argument).
export function mix(...transforms) {
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
