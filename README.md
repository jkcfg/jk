[![Build Status](https://travis-ci.org/jkcfg/jk.svg?branch=master)](https://travis-ci.org/jkcfg/jk)

> `jk` is still in very early stages. The standard library API isn't frozen
> and will most likely change. Feedback is most definitely welcomed!

# Configuration as code

We believe in general purpose languages for configuration. They offer mature
tooling, great runtimes, a well established ecosystem and many learning
resources. We settled for Javascript and built a runtime tailored for
configuration.

# jk

`jk` is a data templating tool designed to help writing structured
configuration. A good way to start with `jk` is to read our [introduction
tutorial][quick-start].

`jk` itself is a Javascript runtime written in Go and embedding [v8][v8].

While a general purpose language is great, configuration code can be made
more maintainable by restricting what it can do. A nice property we can offer
is being "hermetic" that we define with: if you clone a git repository and
execute a `jk` script, the result should be the same on any machine. To give
concrete examples, this means the `jk` standard library doesn't support
environment variables nor has any networking capability.

On top of the `jk` runtime, we are building libraries to help people write
configuration.

- [mixins][mixins]: build and compose configuration objects
- [kubernetes][kubernetes]: build Kubernetes objects

Examples:

- [A Kubernetes deployment written in Typescript][guestbook-ts]
- [Mutating objects with mixins][mixins-example]
- [Kustomize-like behavior][kustomize]

[v8]: https://v8.dev/
[quick-start]: https://jkcfg.github.io/#/documentation/quick-start
[mixins]: https://github.com/jkcfg/mixins
[kubernetes]: https://github.com/jkcfg/kubernetes
[guestbook-ts]: https://github.com/jkcfg/kubernetes/blob/master/examples/guestbook-ts/guestbook.ts
[mixins-example]: https://github.com/jkcfg/mixins/blob/master/examples/mix-simple/namespace.js
[kustomize]: https://github.com/jkcfg/kubernetes/tree/master/examples/overlay
