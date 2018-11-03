#!/bin/sh -e

set -x
container=$1
dest=$2

copy() {
    mkdir -p $dest/`dirname $1`
    docker cp $container:$1 $dest/$1
}

copy /usr/lib/pkgconfig/v8.pc
copy /usr/include/v8
copy /usr/lib/v8
