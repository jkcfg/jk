# jk

[![Build Status](https://travis-ci.org/jkcfg/jk.svg?branch=master)](https://travis-ci.org/jkcfg/jk)

> `jk` is still in very early stages. The standard library API isn't frozen
> and will most likely change. Feedback is most definitely welcome!

## jk - configuration as code

`jk` is a data templating tool designed to help writing structured
configuration files.

The main idea behind `jk` is to use a general purpose language for this task.
They offer mature tooling, great runtimes, a well established ecosystem and
many learning resources. `jk` uses Javascript and a runtime tailored for
configuration.

## Quick start

A good way to start with `jk` is to read our [introduction
tutorial][quick-start].

## More complex examples

- [A Kubernetes deployment written in Typescript][guestbook-ts]
- [Mutating objects with mixins][mixins-example]
- [Kustomize-like behavior][kustomize]

## Architecture & design

### v8

`jk` itself is a Javascript runtime written in Go and embedding [v8][v8]. It
uses Ryan Dahl's [v8worker2][v8worker2] to embed v8 and
[flatbuffers][flatbuffers] for the v8 ⟷ Go communication.

### Hermeticity

While a general purpose language is great, configuration code can be made
more maintainable by restricting what it can do. A nice property `jk` has to
offer is being *hermetic*: if you clone a git repository
and execute a `jk` script, the resulting files should be the same on any
machine. To give concrete examples, this means the `jk` standard library
doesn't support environment variables nor has any networking capability.

### Library support

`jk` provides an unopinionated data templating layer. On top of the `jk`
runtime, libraries provide APIs for users to write configuration.

- [mixins][mixins]: build and compose configuration objects
- [kubernetes][kubernetes]: build Kubernetes objects

## Roadmap

This project is still in early stages but future (exciting!) plans include:

- Reach the state of having Kubernetes examples working and well documented.
- Work on hermeticity. (eg. [#110][issue110], [#44][issue44], [topic/hermeticity][issueHermeticity]).
- Native typescript support ([#54][issue54]).
- HCL support ([#94][issue94]).

[v8]: https://v8.dev/
[quick-start]: https://jkcfg.github.io/#/documentation/quick-start
[mixins]: https://github.com/jkcfg/mixins
[kubernetes]: https://github.com/jkcfg/kubernetes
[guestbook-ts]: https://github.com/jkcfg/kubernetes/blob/master/examples/guestbook-ts
[mixins-example]: https://github.com/jkcfg/mixins/blob/master/examples/mix-simple/namespace.js
[kustomize]: https://github.com/jkcfg/kubernetes/tree/master/examples/overlay
[v8worker2]: https://github.com/ry/v8worker2
[flatbuffers]: https://github.com/google/flatbuffers

[issue44]: https://github.com/jkcfg/jk/issues/44
[issue54]: https://github.com/jkcfg/jk/issues/54
[issue94]: https://github.com/jkcfg/jk/issues/94
[issue110]: https://github.com/jkcfg/jk/issues/110
[issueHermeticity]: https://github.com/jkcfg/jk/issues?q=is%3Aissue+is%3Aopen+label%3Atopic%2Fhermeticity
