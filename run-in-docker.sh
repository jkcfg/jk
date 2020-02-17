#!/bin/bash

pkg=github.com/jkcfg/jk
docker run -v "$(pwd)":/go/src/$pkg jkcfg/build "$@"
