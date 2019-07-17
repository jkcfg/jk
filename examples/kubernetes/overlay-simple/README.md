# Overlays with @jkcfg/kubernetes

This example shows how to compose configuration in the manner of
[Kustomize](https://kustomize.io/), using the
`@jkcfg/kubernetes/overlay` module.

## How to run the example

In this directory,

```console
npm ci
jk generate --stdout ./index.js
```

## What is happening?

This is the simplest way to run an overlay -- it invokes
`generateKustomization` on the current directory, which loads the
kustomization.yaml file there, and interprets it much as Kustomizae
would.
