# Using overlays to compose configurations

This example uses the more generic form of `@jkcfg/kubernetes/overlay`
to _compose_ configurations.

## Running the example

In this directory,

```console
npm ci
jk generate --stdout ./index.js
```

To generate resources for a specific "env", supply a parameter:

```console
jk generate --stdout -p env=staging ./index.js
```

## How does this work?

For the most part, `overlay` mimics Kustomize, letting you load
resources from files, assemble secrets and configMaps, compose via
bases (`kustomization.yaml` files), and transform resources in certain
pre-defined ways, e.g., with `commonLabels`.

But it also provides generic means of composing and transforming, to
cover cases where you want to supply resources you obtained by other
means. You can mix and match the kustomize-like bits with the generic
bits. For example, this lets you instantiate a chart from a module
(using `@jkcfg/kubernetes/chart` say), combine the result with some
files, then patch the results and give them a common namespace and so
on.

This is important because it means you can reuse configurations,
including those from packages.

In the example, most of the resources come from a procedure in
`./submodule`. The submodule generates them by loading all `.yaml`
files in the directory, and running them through the
`@jkcfg/kubernetes/short` expander -- but all that's _required_ for
the overlay is that it returns an array, or a Promise of an array,
containing objects.

Back in the overlay, a namespace for all the resources to live in is
prepended to the submodule resources. Here we don't need a Promise,
just the array is fine.

```js
  const nsResource = new core.v1.Namespace(ns, {});

  return overlay('.', {
    // ...
    generatedResources: [[nsResource], submoduleResources()],
    // ...
  });

```

A touch of kustomization is done with `namespace` and `commonLabels`
-- these work on the generated resources just like they would on
resources from files or bases.

Lastly, the overlay object includes a transformation that adds a
sidecar to any deployments. This is the full overlay object:

```js
  return overlay('.', {
    namespace: ns,
    commonLabels: { env },
    generatedResources: [[nsResource], submoduleResources()],
    transformations: [addSidecar],
  });
```
