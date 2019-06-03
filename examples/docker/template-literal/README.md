# `docker.js`

An example generating a "best practice" Dockerfile using a [template
literal][js-template-literal]. Compared to the [simple
version](../template-literal-simple), the `Dockerfile` specifies a
non-root user to run the application as.

[js-template-literal]: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals

Run this example with:

```console
$ jk generate -v dockerfile.js
wrote Dockerfile
```
