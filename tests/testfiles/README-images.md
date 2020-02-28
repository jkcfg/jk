The `.tar` file(s) in here are container images as tarballs, to be
imported into a temporary registry so they can be used with the flag
`--lib` (in large part, to test whether that flag works).

You can create a tarball from a directory, using Docker, in the
following way (adapted from [Docker's
documentation](https://docs.docker.com/develop/develop-images/baseimages/)). For
a directory `$lib` to be used via the image name `$image`,

```sh
$ tar -c $lib | docker import - $image
$ docker save -o testfiles/$lib.tar $image
```

Assuming the lib provides a module imported as `$lib/foo`, there
should be a file `$lib/foo.js` or `$lib/foo/index.js`. `$lib` itself
can be a path with more than one directory.

The image name when uploaded will be `$lib:v1`, arbitrarily, so should
be used like this:

    jk run --lib $lib:v1 ...

## barlib:v1

barlib is prepared specially so that it contains whiteout files
(including an opaque whiteout).

The layers are like this (from bottom to top):

```
# Adds foo.js and bar.js
ff940e87e638b27658c85f5eea099f2845b22174f4930b0d23dc66f46166a07c
  bar.js
  foo.js

# Adds baz/foo/index.js and whites out everything else in baz/
b3f509bb9ebb848a43fa4753a1ce8a4edcf47a6a8f00f7ef169db2ceee378fa0
  baz/
    .wh..wh..opq
    foo/
      index.js

# Add baz/index.js, which should be visible
41aeafa2014c11a98517a72c5f62bb52565060682b444d773243d8b4b05b4045
  baz/
    index.js

# Adds bar.js and foo.js again (no reason)
01a81ce7b8b6713893656bf8acf4ed5ad454e5d0da599daa2d1e7acf0bedf68f
  bar.js
  foo.js

# Adds baz/index.js and whites out everything else in baz/
# (hiding baz/foo/index.js)
4adec11db178e60ff80fdb392ae4c2c7e15dfdb90cf113705c580caa05ff22c4
  baz/
    .wh..wh..opq
    index.js

# Whites out bar.js with .wh.bar.js
67ab3fd25033b0983417d7094762e506c418bf1a834c5f4410e9e68daaa0ea84
  .wh.bar.js

# Overlays foo.js with another file (default export 'baz')
1a18b3cf1e651323c4d059150b976eca9759bbfa8af1b39c7abf784abb005e19
  foo.js
```

The resulting filesystem should look like this:

```
foo.js # exports 'baz'
# bar is whited-out
baz/
  index.js
  # foo is opaque-whited-out
```
