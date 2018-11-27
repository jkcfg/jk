#!/bin/bash

pkg=github.com/justkidding-config/jk
docker run -v "$(pwd)":/go/src/$pkg quay.io/justkidding/build "$@"
