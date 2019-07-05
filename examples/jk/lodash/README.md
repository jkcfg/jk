# Using `lodash` with `jk`

`jk` can import an npm module if:

- it uses [es6 modules](https://hacks.mozilla.org/2015/08/es6-in-depth-modules/),
- it doesn't use nodejs APIs.

There are many of those npm modules!

To use `lodash` we need a es6 version of that package. Fortunately, such a
module exists: `lodash-es`.

Run the lodash example with:

```
$ npm install lodash-es
$ jk run lodash.js
```
