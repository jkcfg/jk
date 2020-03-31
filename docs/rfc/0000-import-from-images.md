# Using container images for packaging (experimental)

## Summary

This RFC conjectures that container images are a useful form of
packaging (for discussion, see the rest of the text), and proposes:

 - adding a flag for specifying a container image to put on the module
   search path
 - automatically downloading and caching images that appear on the
   module search path

The main implementation changes are:

 - allowing module resolvers to access container images
 - an image download and caching mechanism

The latter is orthogonal to its use for fetching modules, though this
is the only use for it so far.

This RFC was developed in conjunction with an implementation in [PR
315][#315].

## Example

```sh
$ jk validate --lib jkcfg/kubernetes:0.6.0 ./lint.js *.yaml
```

This adds the filesystem in the image `jkcfg/kubernetes:0.2.1` to the
module search path, then runs the validation function from `./lint.js`
(which is for the sake of the example, assumed to import modules from
the image) over the files given by the glob pattern.

`jk` will check if the image is already cached locally, and if not,
download it before running.

## Motivation

Presently, `jk` will resolve imports by looking in the filesystem
starting at the directory containing the script being run. This works
in sympathy with NPM, so long as you have a package.json describing
your dependencies.

However, there are some reasons to be dissatisfied with NPM, and some
reasons to explore using container images for packaging and
distribution. (See Alternatives for some more discussion of NPM).

The main kind of use case this addresses is how to arrange for
dependencies to be available along with code. For example, when using
`jk` with [Flux's manifest generation][flux-manifest], it's not enough
for a generation script to be present in the git repository, since it
will almost always need some libraries to be there as well.

There are a few ways to arrange this:

 - vendor the libraries (i.e., add them all to the git repository)
 - bake NPM into the Flux container, and run it before running the
   script
 - copy files from an initContainer into a volume in the module search
   path
 - create a custom image that includes fluxd, jk, and the libraries
 - ... and no doubt, variations on the above.

All require some work outside of the code and invocation of `jk`
itself, and none recommend themselves as elegant or convenient.

The RFC proposes a way to make fetching dependencies convenient,
without sacrificing repeatability, and introduces a low-stakes way to
trial using container images for packaging.

## Design

### User interface

    jk run --lib jkcfg/kubernetes:0.2.1 -I ../mycharts generate.js

`--lib` is put forward here as small and guessable. `--image` is an
alternative, but that doesn't indicate it's a _library_ rather than
the thing to execute.

### Downloading and caching images

When a library image is required for excuting a `jk` script, these
things shall happen:

 1. The image ref is resolved to a path within the cache directory

    1.2. whether the ref is a tag or a digest, the path is expected to
         be a symlink to a blob which is the manifest. Note that the
         digest ref for an image is not the digest of the manifest
         blob (I don't know what it is, maybe the digest of the
         gzipped file?)

 2. If the file does not exist, it is resolved using OCI distribution,
    and written into the cache (as per the expectations in the
    previous step)

 3. The file at the path, or linked, is a manifest (MIME type
    application/vnd.oci.image.manifest.v1+json); if necessary, it's
    converted to that from a Docker image manifest

    3.1. If the fetched file is a manifest list, the appropriate
         manifest from it is fetched in turn.

 4. For each layer in the image (we don't care about the config,
    that's for running the image), if it is missing then fetch it,
    verify its digest, and put it in the cache blob store

 5. A filesystem based on the image layers is added to the module
    search path.

### Layout of images

Container images include a whole filesystem. Where should `jk` look
for modules?

 * we will want image layers to compose so they can be remuxed, so it
   makes sense to put modules in the same place every time
 * e.g., just chuck everything under `/jk/modules/`, so if you jam a
   bunch of layers together in an image, it'll just look like a
   directory with all the libraries under `/jk/modules/`.

### Using images in the module search path

In two steps:

 - a. adapt the module resolvers to allow filesystems represented by
   http.FileSystem (or similar)
 - b. write a union filesystem that understands the diff layers of an OCI
   image

Since modules may be in some abstraction of a filesystem, access to
resources will need to go through the abstraction too, as will any
other built-in code that loads files (like schema validation).

The work for a.) was done in [PR 307](#307). The remainder, i.e., b.)
is done in [PR 315][#315].

#### Building images

Although this RFC is not about tooling for building images, it's worth
mentioning here how one would go about it.

Using a Dockerfile, you can make an archive of your current
node_modules contents:

```
$ cat > Dockerfile <<EOF
FROM scratch
WORKDIR /jk/modules
# NB COPY copies the _contents_ of directories
COPY node_modules ./
EOF
$ docker build -t localhost:5000/mylib:0.1-pre .
$ docker push localhost:5000/mylib:0.1-pre
```

It is also desirable to be able to combine libraries -- this is an
assumed benefit of using image layers, that you can (in some
circumstances) combine them, even without having the layers
present. That is outside the capability of `docker build`, and may
make a good addition to `jk`, later.

[Buildpacks](https://buildpacks.io/) are another possibility.

## Backward compatibility

Existing uses of `jk` won't need altering, though might benefit from
being adapted to use images.

## Drawbacks and limitations

_Does it require more engineering work from users (for an equivalent
result)?_

Using images for libraries is entirely optional. The premise is that
it will make some things easier, though there is work needed to take
advantage of it.

It is fairly easy to create an image that you can use with `jk --lib`,
and automation can be added to build and push them to an image
registry (`@jkcfg/kubernetes` already does this). In general, the
burden falls on library developers (i.e., mainly the `jk` authors)
rather than users.

_Are there performance concerns?_

This feature is strictly user-pays -- `jk` works exactly as before, if
you don't use images.

There is an initial delay when using an image that is not in the
cache, so running time is less predictable. Extra tooling for pulling
images ahead of time might be useful to offset this.

_Will it close off other possibilities?_

In the sense that there would now be a "packaging" solution built into
`jk`, and it would be odd to have more than one, this is a step down a
particular path, yes.

_Does it add significant complexity to the runtime or standard
library?_

The implementation of downloading, caching, and using their
filesystems is a few hundred lines of code (see [PR 315][#315]). The
effect on existing code is relatively small, since it just involves
adding resolvers to the module search path.

_Does it make understanding `jk` harder?_

I would not think so. Some confusion might arise over referring to
images and referring to modules, since these will tend to have similar
names:

    jk run --lib jkcfg/kubernetes:0.2.1 -m @jkcfg/kubernetes/validate ...

Perhaps there's room to finesse the user interface (e.g., to map
modules to images in a dot file), to help with that confusion.

## Alternatives

### Stick with NPM

An obvious question is "Why not just bless NPM as the way to
distribute libraries for jk?". That could entail downloading NPM
modules automatically, too.

One downside of using NPM is that it's understandably weighted towards
Node.JS. In practice it will work for resolving and downloading
packages for `jk`, but there's a bunch of cognitive baggage, like
irrelevant fields in `package.json`.

Another downside is that NPM packages are not architecture- or
OS-specific -- like Node.JS, they are platform-independent. For
plugins, and perhaps other purposes, `jk` will occasionally need to
ship multiple binaries. While it's possible to dispatch within the
runtime given a package with all the binaries, it'd be nicer to just
download a platform-specific package.

In its favour, NPM and its quirks are well known by now, and you can
just use the fields you need ('name' and 'dependencies' more or less),
and it'll work for most packages.

Using container images for packaging introduces some interesting
possibilities. Since an image is a set of layers with some metadata,
it's possible to construct a single image with all dependencies which
nonetheless shares structure (layers) with other images, and can
thereby be distributed efficiently.

It's also possible to make platform-specific images (that share
layers) to save people fetching redundant files.

However, there's no dependency resolution machinery, and it may be
fiddly to co-opt that of NPM (if it were considered suitable).

So it comes down to: are the possibilities of using images
sufficiently interesting to try using them in the way detailed here,
with the prospect of building on that later. (Answer: Sure!)

### Specify images in imports rather than on the command line

In this RFC, an import in JavaScript is given as a symbolic name,
which is resolved to a location depending on the command-line flags
given to `jk` controlling the module search path.

Another way is to imply the location in the import statement itself
(like golang does). For example, instead of

```
import { chart } from '@jkcfg/kubernetes/helm';
```

use

```
import { chart } from 'oci:jkcfg/kubernetes/helm:0.5.2';
```

which gives the provenance of the imported module as being the image
`jkcfg/kubernetes/helm:0.5.2` (assuming that could be resolved to an
image repo somewhere). An advantage of this is that it make statically
determining the dependencies easier: you can in principle just take
the code, download the images involved, look at _those_ modules, etc.

There are reasons you might want to keep the indirection, though; most
dependency resolution tooling keeps the specifics in a separate file
(go.mod/go.sum; various lockfiles), and anything we do in `jk` is
likely to do the same. This RFC does not rule out being more explicit
about provenance in imports later.

### Use "OCI distribution" and own format

Helm 3 has experimental support for pushing charts (i.e., tarballs) to
an OCI registry. Here is a library for pushing arbitrary artifacts to
an OCI registry: https://github.com/deislabs/oras

On the registry side, it's not ready for primetime yet: it's only
supported by ACR and registry:2.

But the main argument against this is that we are interested in some
properties of container images specifically, not just of
registries. Using another format would mean either reinventing things
or losing those properties.

## Unresolved questions

**tags vs checksums/digests**

It's possible for the image that a ref points at to change. What
should we do in that circumstance (will we even notice? once it's
downloaded, we probably won't check again). This could be punted by
letting people update the cached images explicitly, and/or reporting
when the digest has changed. A `go.sum`-style file might be used to
make things repeatable. Or a digest given in the image name (but then
.. so much text .. maybe getting the search path from a file would be
useful).

[#315]: https://github.com/jkcfg/jk/pull/315
[#307]: https://github.com/jkcfg/jk/pull/307
[flux-manifest]: https://docs.fluxcd.io/en/stable/references/fluxyaml-config-files.html
