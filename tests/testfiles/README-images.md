The `.tar` file(s) in here are container images as tarballs, to be
imported into a temporary registry so they can be used with the flag
`--lib` (in large part, to test whether that flag works).

You can create a tarball from a directory, using Docker, in the
following way (adapted from [Docker's
documentation](https://docs.docker.com/develop/develop-images/baseimages/)). For
a filesystem under `$src` to be used via the image name `$image`,

```sh
$ tar -C $src -c jk | docker import - $image
$ docker save -o $image.tar $image
```

`jk` modules are expected to be under the path `/jk/modules/` in the
image filesystem. Assuming the lib provides a module imported as
`path/to/foo`, there should be a file `$src/jk/module/path/to/foo.js`
or `$src/jk/modules/path/to/foo/index.js`.

The image name when _uploaded_ will be the base name of the tarfile,
with a tag `v1`; i.e., `$image:v1`. The registry URL will be in the
environment variable $REGISTRY, so in a test you can refer to, for
example,

    jk run --lib ${REGISTRY}/$image:v1 ...

## foolib:v1

The Makefile uses the above recipe to create `foolib.tar` from
`./src/foolib/` as the source path.

## barlib:v1

barlib is prepared specially so that it contains whiteout files
(including an opaque whiteout). This needs a bit of arranging:

 1. [`Dockerfile.base-foo`](./src/Dockerfile.base-foo) copies `foo.js`
    and `bar.js`, and creates the directory `baz/` implicitly by
    copying to `baz/baz.js`.
 2. [`Dockerfile.alpine-foo`](./src/Dockerfile.alpine-foo) copies
    foo.js and baz.js (but _not_ baz/baz.js), and is based on
    alpine:3.9 so we have `rm` available;
 3. [`Dockerfile.alpine-bar`](./src/Dockerfile.alpine-bar) is based on
    `alpine-foo`, so it has `foo.js` and `bar.js`. It removes `bar.js`,
    creating a whiteout for that file, and replaces `foo.js`. It _also_
    creates a directory `baz/`, which, since it's not present in
    `alpine-foo`, makes an opaque whiteout;
 4. `alpine-bar` is rebased on `base-foo`, removing the alpine layers
    but leaving the opaque whiteout from step 3.

The resulting filesystem looks like this:

```
foo.js # exports 'baz', since replaced
# bar.js is whited-out, since removed
baz/
  index.js
  # baz.js is opaque-whited-out
```
