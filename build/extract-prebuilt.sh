#!/bin/sh -e

set -x

image=quay.io/justkidding/build
container=jk-extractor

docker run --rm -d --name $container $image sleep 3600

dest=$1
[ -z "$dest" ] && dest=.

copy() {
    docker cp $container:$1 $2
}

copy /usr/local/lib/pkgconfig/v8.pc $dest/
rm -rf $dest/include
copy /usr/local/include/v8 $dest/include
copy /usr/local/lib/v8/libv8_monolith.a $dest/
copy /usr/local/bin/flatc $dest/

docker kill $container
