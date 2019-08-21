# Releasing `jk`

Steps to produce a new `jk` version:

1. Checkout the master branch and make sure it reflects the latest `origin/master`:

   ```console
   $ git checkout master
   $ git pull --rebase
   ```

1. Bump the version in `std/package{-lock}.json`, commit and push the result:

   ```console
   $ vim std/package.json # edit the 'version' field
   $ vim std/package-lock.json # edit the 'version' field
   $ git commit -a -m "build: Bump std package version to x.y.z"
   $ git push
   ```

1. Create and push the new tag:

    ```console
    $ git tag -a x.y.z -m x.y.z
    $ git push --tags
    ```

1. Wait for CI to successfully push the release binaries and npm module.

1. Redact the release changelog on github with the list of new features, API
changes and bug fixes.

1. Freeze the API refrence documentation for this version:

   ```console
   # TODO(damien): automate this in run-release.sh
   $ git clone git@github.com:jkcfg/jkcfg.github.io.git
   $ cd jkcfg.github.io
   $ cp -r reference/std/latest reference/std/x.y.z
   $ git add reference/std/x.y.z
   $ git commit -m '@jkcfg/std: Freeze API reference for x.y.z'
   $ git push
   ```

1. Add a new link to the API reference on the [website][website] in
`static/docs/std-api-reference.md`.

1. Bump the latest version in `src/params.json` located in the [website
repository][website] and push the result.


[website]: https://github.com/jkcfg/website
