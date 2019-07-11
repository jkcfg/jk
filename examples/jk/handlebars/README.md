# Using `handlebars` with `jk`

`jk` can import an npm module if:

- it uses [es6 modules](https://hacks.mozilla.org/2015/08/es6-in-depth-modules/),
- it doesn't use nodejs APIs.

There are many of those npm modules!

To use `handlebars` we need to import its top level es6 module:
`handlebars/lib/handlebars`.

Run the handlebars example with:

```
$ npm install handlebars
$ jk generate -v handlebars.js
wrote index.html
```
