# Examples of using `jk` to generate Kubernetes manifests

This directory contains examples of using `jk`, with the library
`@jkcfg/kubernetes`, to generate Kubernetes configuration.

## Running an example

Each example is in its own directory, with a `package.json`, and
`index.js`. To run an example, cd to the example directory and install
the dependencies:

    cd <example>/
    npm ci

The examples are all designed to work with `jk generate`:

    jk generate --stdout ./index.js

will run the example and print the result (usually as a YAML stream)
to stdout. The example packages are set up so that `npm run generate`
will run the above.
