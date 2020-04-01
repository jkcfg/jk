# Test files for image cache

This directory contains

 - a preprepared image cache (`./dotcache`)
 - tarballs of images for uploading to a temporary registry, to test
   downloading file to the cache


The latter are prepared in different ways.

`helloworld.tar` is simply the Docker hello-world image, pulled then
saved:

    docker pull hello-world:linux
    docker save -o helloworld.tar hello-world:linux

`symlink.tar` is prepared by importing a tar file:

    tar -c symlink | docker import - symlink:v1
    docker save -o symlink.tar symlink:v1

