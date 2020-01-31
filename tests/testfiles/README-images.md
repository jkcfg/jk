The `.tar` file(s) in here are container images as tarballs, to be
imported into a temporary registry so they can be used with the flag
`--lib` (in large part, to test whether that flag works).

The tarballs are created from a directory, using Docker, in the
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
