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

function mergeFunc(rule, key, defaultFunc) {
  if (rule === undefined) {
    return defaultFunc;
  }

  const f = rule[key];
  if (f === undefined) {
    return defaultFunc;
  }

  if (typeof f !== 'function') {
    throw new Error(`merge: expected a function in the rules objects but found a ${typeof f}`);
  }

  return f;
}

function objectMerge2(a, b, rules) {
  const r = {};

  Object.assign(r, a);
  for (const [key, value] of Object.entries(b)) {
    r[key] = mergeFunc(rules, key, mergeFull)(a[key], value);
  }
  return r;
}

/**
 * Merge strategy deep merging objects.
 *
 * @param rules optional set of merging rules.
 *
 * `deep` will deep merge objects. This is the default merging strategy of
 * objects. It's possible to provide a set of rules to override the merge
 * strategy for some properties. See [[mergeFull]].
 */
export function deep(rules) {
  return (a, b) => objectMerge2(a, b, rules);
}

function arrayMergeWithKey(a, b, mergeKey, rules) {
  const r = Array.from(a);
  const toAppend = [];

  for (const value of b) {
    const i = a.findIndex(o => o[mergeKey] === value[mergeKey]);
    if (i === -1) {
      // Object doesn't exist in a, save it in the list of objects to append.
      toAppend.push(value);
      continue;
    }
    r[i] = objectMerge2(a[i], value, rules);
  }

  Array.prototype.push.apply(r, toAppend);
  return r;
}

/**
 * Merge strategy for arrays of objects, deep merging objects having the same
 * `mergeKey`.
 *
 * @param mergeKey key used to identify the same object.
 * @param rules optional set of rules to merge each object.
 *
 * **Example**:
 *
 * ```js
 * import { mergeFull, deep, deepWithKey } from '@jkcfg/std/merge';
 *
 * const pod = {
 *   spec: {
 *     containers: [{
 *       name: 'my-app',
 *       image: 'busybox',
 *       command: ['sh', '-c', 'echo Hello Kubernetes!'],
 *     },{
 *       name: 'sidecar',
 *       image: 'sidecar:v1',
 *     }],
 *   },
 * };
 *
 * const sidecarImage = {
 *   spec: {
 *     containers: [{
 *       name: 'sidecar',
 *       image: 'sidecar:v2',
 *     }],
 *   },
 * };
 *
 * mergeFull(pod, sidecarImage, {
 *   spec: deep({
 *     containers: deepWithKey('name'),
 *   }),
 * });
 * ```
 *
 * Will result to:
 *
 * ```js
 * {
 *   spec: {
 *     containers: [
 *       {
 *         command: [
 *           'sh',
 *           '-c',
 *           'echo Hello Kubernetes!',
 *         ],
 *         image: 'busybox',
 *         name: 'my-app',
 *       },
 *       {
 *         image: 'sidecar:v2',
 *         name: 'sidecar',
 *       },
 *     ],
 *   },
 * }
 * ```
 */
export function deepWithKey(mergeKey, rules) {
  return (a, b) => arrayMergeWithKey(a, b, mergeKey, rules);
}

/**
 * Merges `b` into `a` with optional merging rule(s).
 *
 * @param a Base value.
 * @param b Merge value.
 * @param rule Set of merge rules.
 *
 * `mergeFull` will recursively merge two values `a` and `b`. By default:
 *
 * - if `a` and `b` are primitive types, `b` is the result of the merge.
 * - if `a` and `b` are arrays, `b` is the result of the merge.
 * - if `a` and `b` are objects, every own property is merged with this very
 * set of default rules.
 * - the process is recursive, effectively deep merging objects.
 *
 * if `a` and `b` have different types, `mergeFull` will throw an error.
 *
 * **Examples**:
 *
 * Merge primitive values with the default rules:
 *
 * ```js
 * mergeFull(1, 2);
 *
 * > 2
 * ```
 *
 * Merge objects with the default rules:
 *
 * ```js
 * const a = {
 *   k0: 1,
 *   o: {
 *     o0: 'a string',
 *   },
 * };
 *
 * let b = {
 *   k0: 2,
 *   k1: true,
 *   o: {
 *     o0: 'another string',
 *   },
 * }
 *
 * mergeFull(a, b);
 *
 * >
 * {
 *   k0: 2,
 *   k1: true,
 *   o: {
 *     o0: 'another string',
 *   }
 * }
 * ```
 *
 * **Merge strategies**
 *
 * It's possible to override the default merging rules by specifying a merge
 * strategy, a function that will compute the result of the merge.
 *
 * For primitive values and arrays, the third argument of `mergeFull` is a
 * function:
 *
 * ```js
 * const add = (a, b) => a + b;
 * mergeFull(1, 2, add);
 *
 * > 3
 * ```
 *
 * For objects, each own property can be merged with different strategies. The
 * third argument of `mergeFull` is an object associating properties with merge
 * functions.
 *
 *
 * ```js
 * // merge a and b, adding the values of the `k0` property.
 * mergeFull(a, b, { k0: add });
 *
 * >
 * {
 *   k0: 3,
 *   k1: true,
 *   o: {
 *     o0: 'another string',
 *   }
 * }
 * ```
 */
export function mergeFull(a, b, rule) {
  const [typeA, typeB] = [typeof a, typeof b];

  if (a === undefined) {
    return b;
  }

  if (typeA !== typeB) {
    throw new Error(`merge cannot combine values of types ${typeA} and ${typeB}`);
  }

  // Primitive types and arrays default to being replaced.
  if (Array.isArray(a) || typeA !== 'object') {
    if (typeof rule === 'function') {
      return rule(a, b);
    }
    return b;
  }

  // Objects.
  return objectMerge2(a, b, rule);
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
