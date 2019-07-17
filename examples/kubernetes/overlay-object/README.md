# Overlay constructed in code

The `overlay-simple` example ran a kustomization file by referring to
a kustomization file in its directory `'.'`. This being JavaScript, we
should be able to construct an overlay programmatically.

This example uses the same kustomization file, and base manifest, as
before, but overlays changes on top.

## Running the example

In this directory,

```
npm ci
jk generate --stdout ./index.js
```

## What is happening?

To use the original kustomization file, a `base` is included:

```js
const kustom = {
  bases: ['.'],
  // ... see below for the rest
}
```

If we stopped there, it would be the same kustomization as
`overlay-simple` (because it's the same files). In this example
though, there's further work to be done. Here's the full kustomization
object:

```js
const kustom = {
  bases: ['.'],
  // This adds a resource loaded from the file mentioned
  resources: ['service.yaml'],

  // This adds a label to all resources
  commonLabels: {
    team: 'strange',
  },

  // This patches the deployment from the original kustomization,
  // so that it will be selected by the service added above
  patches: ['service-selector.json'],

  // This sets the namespace for all the resources to a param supplied
  when the script is run (or a default, if not supplied)
  namespace: param.String('namespace', 'default'),
}
```

Most of this is a translation of what you'd see in a kustomization
file -- it's what you'd get if you parsed the file into a JavaScript
object. The last field though, `namespace:`, is given a calculated
value. The calculation is this: if a value was supplied as a parameter
(using the `-p` or `-f` flags with `jk generate`), use that, otherwise
set it to `default`.

You can see what effect the parameter has by running

```console
jk generate --stdout -p namespace=demo ./index.js
```
