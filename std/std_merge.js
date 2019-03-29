
// patch returns a new value that has the fields of `obj`, except
// where overridden by fields in `patchObj`. Entries in common are
// themselves patched. This is similar to `merge` below, but always
// does does a deep merge.
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

// merge transforms `obj` according to the field given in
// `mergeObj`. A field name ending in '+' is "deep merged", that is,
// patched; otherwise, the value of the field is simply assigned into
// the result. Any other fields in `obj` are also assigned in the
// result.
function merge(a, b) {
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

export { patch, merge, mix };
