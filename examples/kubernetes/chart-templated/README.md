# Helm chart analogue using templates

This example shows a "chart" that uses textual templates to generate
the resources. In the `chart-simple` example, the resources were
object literals; in this example, they are loaded from files in a
directory relative to the chart module.

The chart module (in the directory `chart`) shows how you can package
the ingredients of chart together -- in principle, it could be made
into its own, self-contained NPM package.

## To run the example

```console
# install the dependencies
npm ci

# run the example
jk generate --stdout ./index.js

# run the example with some parameters set
jk generate --stdout ./index.js -p values.name=demo -values.image.tag=1.1
```

## Explanation

In this example, all the bits to do with the chart are in the
subdirectory `chart`. The index file there uses the module
`@jkcfg/kubernetes/chart/template` to construct manifests given a
directory of template files; and, loads the value defaults from a file
(also relative to the module) `default.yaml`.

### Templates

In principle, any templating engine could be used, so long as there is
a function for turning a string (or file path) into a function for
instantiating the template. In practice, not every library will be
convenient to use, since `jk` uses ES6 modules, and many libraries
only distribute CommonJS modules.

Happily, `handlebars` distributes ES6 modules, _and_ has a similar
template syntax and usage to `gotpl` (the predominant templating used
for Helm charts). `handlebars/lib/handlebars#compile` is exactly the
function we need to turn a string into a template. In the example
code, `compile` is given to `loadModuleTemplates`, a procedure that
will load all the templates relative to a module. This way, the
templates can be distributed with the module, e.g., in an NPM package.
