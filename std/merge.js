/**
 * @module std/merge
 */

function mergeFunc(rule, key, defaultFunc) {
  const f = rule && rule[key];
  if (f === undefined) {
    return defaultFunc;
  }

  const t = typeof f;
  if (t === 'object' && t !== 'function') {
    return deep(f);
  }
  if (t !== 'function') {
    throw new Error(`merge: expected a function in the rules objects but found a ${t}`);
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

function isObject(o) {
  const t = typeof o;
  return t === 'object' && !Array.isArray(o) && !!o;
}

function assertObject(o, prefix) {
  if (!isObject(o)) {
    throw new Error(`${prefix}: input value is not an object`);
  }
}

function isArray(a) {
  return Array.isArray(a);
}

function assertArray(o, prefix) {
  if (!isArray(o)) {
    throw new Error(`${prefix}: input is not an array`);
  }
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
  return (a, b) => {
    assertObject(a, 'deep');
    assertObject(b, 'deep');
    return objectMerge2(a, b, rules);
  };
}

/**
 * Merge strategy merging two values by selecting the first value.
 *
 * **Example**:
 *
 * ```js
 * let a = {
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
 * };
 *
 * mergeFull(a, b, { o: first() });
 * ```
 *
 * Will give the result:
 *
 * ```js
 * {
 *   k0: 2,
 *   k1: true,
 *   o: {
 *     o0: 'a string',
 *   },
 * }
 * ```
 */
export function first() {
  return (a, _) => a;
}

/**
 * Merge strategy merging two values by selecting the second value.
 *
 * **Example**:
 *
 * ```js
 * let a = {
 *   k0: 1,
 *   o: {
 *     o0: 'a string',
 *     o1: 'this will go away!',
 *   },
 * };
 *
 * let b = {
 *   k0: 2,
 *   k1: true,
 *   o: {
 *     o0: 'another string',
 *   },
 * };
 *
 * mergeFull(a, b, { o: replace() });
 * ```
 *
 * Will give the result:
 *
 * ```js
 * {
 *   k0: 2,
 *   k1: true,
 *   o: {
 *     o0: 'another string',
 *   },
 * }
 * ```
 */
export function replace() {
  return (_, b) => b;
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
 *   spec: {
 *     containers: deepWithKey('name'),
 *   },
 * });
 * ```
 *
 * Will give the result:
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
  return (a, b) => {
    assertArray(a, 'deepWithKey');
    assertArray(b, 'deepWithKey');
    return arrayMergeWithKey(a, b, mergeKey, rules);
  };
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
