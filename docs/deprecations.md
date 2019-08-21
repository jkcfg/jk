# List of std library deprecations

## Deprecated in 0.3.x (will be removed in 0.4.0)

### `generate` file property is now called path

The `file` property in the array of objects consumed by `jk generate` is
deprecated in favour of the `path` property. Both names still work.

**Deprecated**

```
const object = {
  message: 'success',
};

export default [
  { file: 'object0.yaml', value: object },
];
````

**Use**:
```
export default [
  { path: 'object0.yaml', value: object },
];
```

## Deprecated in 0.2.x (will be removed in 0.3.0)

### merge, patch and mix std functions

*Deprecated in 0.2.10*

The `merge` and `patch` functions of the `@jkcfg/std/merge` module have been
deprecated in favour of the more general `mergeFull` function.

The `mix` function has been deprecated to be redefined a bit later. Composing
merge operations is definitely useful but we'd like to think a bit more about
it.

- `merge` and `patch` will be removed in `0.3.0`.
- `mergeFull` will be renamed to `merge` in `0.3.0`.
- `mix` will be removed in `0.3.0`.

### std import

*Deprecated in 0.2.5*

We have decided the std library should be imported from `'@jkcfg/std'`:

- It lets us write a node.js shim at a later point that replicates the
  stdlib functionality for use in unit tests.
- We publish typescript typings in this package, `tsc` will pick up those
  definitions automatically as well as IDEs.

We also decided to [not use `export default`](https://basarat.gitbooks.io/typescript/docs/tips/defaultIsBad.html).

**Deprecated**:

```js
import std from 'std';
import std from '@jkcfg/std';
```

**Use:**

```js
import * as std from '@jkcfg/std';
```

or for specific functions:

```js
import { log } from '@jkcfg/std';
```

This command can be used to help porting existing code over:

```
find . -name "*.js" -o -name "*.ts" -o -name "*node_modules*" -prune | \
  xargs sed -i -e "s#import std from 'std';#import * as std from '@jkcfg/std';#"
```

### std sub-modules

*Deprecated in 0.2.5*

`jk` can now use fine grained modules in its standard library. We have split out:

- '@jkcfg/std/param'
- '@jkcfg/std/fs'

**Deprecated**

```js
import { param, dir, info }  as std from '@jkcfg/std';
```

**Use:**

```js
import * as param from '@jkcfg/std/param';
import { dir, info } from '@jkcfg/std/fs';
```
