#!/bin/bash

pkg=github.com/jkcfg/jk
docker run -v "$HOME/.npm":/go/src/$pkg/.npm -v "$(pwd)":/go/src/$pkg jkcfg/build "$@"
