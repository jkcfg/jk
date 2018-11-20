#!/bin/bash

pkg=github.com/dlespiau/jk
docker run -v "$(pwd)":/go/src/$pkg quay.io/justkidding/build "$@"
