# List of std library deprecations

## Deprecated in 0.2.x (will be removed in 0.3.0)

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
