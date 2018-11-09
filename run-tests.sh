#!/bin/sh -e

(
  cd std
  npm install
  mkdir -p build
  npm run build
)

# Run std tests
(
  cd std
  npm test
)

# Build jk
export GO111MODULE=on
go install

# Tests, both unit tests and integration tests under /tests
go test -v ./...
