# Plugins

## Summary

We have found a number of use cases where we'd like to extend the
functionality of `jk` but don't want to link the core runtime w specific
libraries.

With a plugin system, we could implement these features outside of the
runtime with and provide simple way to extend `jk` and experiment with ideas.

The proposal is to:

- Use [go-plugin](https://github.com/hashicorp/go-plugin), the plugin system
underpinning terraform.
- Define a first integration point to extend `jk` functionality: `Render`
plugins.
- Implement a helm plugin in `@jkcfg/kubernetes` that renders helm charts and
make its Kubernetes objects available to the `jk` runtime for further
manipulation.

## Example

1. A [helm](https://helm.sh/) `Renderer` plugin would render a helm chart
from values specified in `js` and return an array of Kubernetes objects.

1. A Dockerfile `Validator` plugin could validate Dockerfiles that we
write.

1. An [Open Policy Agent](https://www.openpolicyagent.org/) `Validator`
plugin would ensure a set of configuration files passes a policy written in
[Rego](https://www.openpolicyagent.org/docs/latest/policy-language/)

## Motivation

Let's take the helm `Renderer` example: the community is producing some
quality charts users would like to be able to reuse.

- It is counterproductive to start translating complex helm charts in `js`
instead of reusing them.
- It is often necessary to modify the Helm Chart Kubernetes objects in ways
the original authors haven't thought of. Importing them in `jk` allows just
that.

## Design

### General Flow

A call to a plugin a really an RPC call issued to a plugin server process.
`jk` is responsible for the life cycle of those plugin processes.

From an API point of view, it is envisioned that library would provide nice
wrapper objects. For instance the helm chart renderer could look like:

```js
const redis = new k8s.helm.Chart("redis", {
    repo: "stable",
    chart: "redis",
    version: "3.10.0",
    values: {
        usePassword: true,
        rbac: { create: true },
    },
 });

redis.render().then(std.log);
```

The `Chart` object, in this case, would be part of the `@jkcfg/kubernetes` library.

The `render()` function of the `Chart` is implemented with a standard library
RPC call that lands in the go part of `jk`. Then:

- `jk` checks if a helm renderer plugin is already running and spawns a new one if needed
- `jk` waits until the plugin process has started and is ready to accept RPCs.
- `jk` issues a RPC call to the plugin process
- The plugin process returns an answer
- The answer is serialized and sent back to the `js` vm

### The standard library `plugin` function

The core `plugin` RPC call is quite general:

```text
plugin(kind: string, url: string, input: JSON) -> JSON
```

- `kind` is the kind of plugin invoked (eg. `render`). Plugin binaries can
implement more than one kind of plugins.
- `url` identifies a plugin, for instance:
`https://jkcfg.github.io/plugins/helm/0.1.0/plugin.json`. This JSON file is
really plugin metadata (see next section). It is highly recommended, for
reproducibility, to ensure the plugin definition and binaries the URL points
to be immutable to encode a version in the URL.
- `input` and the return value are generic JSON objects that wrapping library
code are responsible for understanding.

### Plugin definition

The plugin URL in the `plugin` call points to a JSON file describing the plugin:

```json
{
  "name": "helm",
  "version": "0.1.0",
  "binaries": {
    "linux-amd64": "https://jkcfg.github.io/plugins/helm/0.1.0/linux-amd64/jk-plugin-helm",
    "darwin-amd64": "https://jkcfg.github.io/plugins/helm/0.1.0/darwin-amd64/jk-plugin-helm"
  }
}
```

### Local plugins

For development purposes it is possible to point the `plugin` RPC call to a
local file by giving a relative path to the JSON plugin definition. The
`binaries` fields can point at plugins present in the `PATH`.

```json
{
  "name": "echo",
  "version": "0.1.0",
  "binaries": {
    "linux-amd64": "jk-plugin-echo",
    "darwin-amd64": "jk-plugin-echo"
  }
}
```

## Backward compatibility

New feature, no backward compatibility concerns.

## Drawbacks and limitations

### Complexity

Plugins do add more complexity to `jk`. Complexity creep is somewhat
unavoidable if we want to support things such as consuming helm charts. One
good thing about plugins is that, at the cost of a (simple) interface between
`jk` and plugins, most of the complexity is delegated to plugin code, not the
core.

We should ensure that plugins don't add any cognitive load on the user. By
wrapping plugin invocation in library objects we can make then mostly
transparent to the user with, maybe, the exception of plugin downloading.

## Hermeticity

Plugins have a high potential to break hermeticity. We should ensure our
plugins are made "in good faith", are self-contained and as deterministic as
possible.

I believe that's an ok price to pay for such extensibility power.

## Alternatives

- A generic `exec` standard library function that would execute any binary in
  the path.

  Problems with that approach:

  1. Executing whatever is in the path doesn't play well with hermeticity.
  Plugins have versions for reproducibility.
  1. Packaging. One goal is to be able to package all the dependencies needed
  for a `jk` script to run. Having the plugin abstraction with metadata
  allows that.

## Unresolved questions

- How to download plugins. I'd like to have very little friction when using
plugins. `jk` should download plugins, somehow. This is linked to a more
general problem of downloading all needed dependencies to run a `jk` script.

- It would be nice to be able to cache dependencies/artifacts plugins needs
beyond the plugin binary itself. For instance, in the case of the helm
renderer, the plugin could cache the chart so subsequent `jk` runs don't need
to hit the network. While the plugin itself could do that caching, it'd be
even nicer if `jk` could help with that: the plugin could ask `jk` to cache
things on its behalf. We'd then be able to gather all dependencies in one
place for `jk` runs that don't hit the network at all.

- Similarly, `jk` could provide plugins a download API so downloading +
caching artifacts UX is consistent across plugins.
