#!/bin/sh -e

echo "==> Building std"
(
  cd std
  npm install
  mkdir -p build
  npm run build
)

echo "==> Running std tests"
(
  cd std
  npm test
)

echo "==> Building jk"
export GO111MODULE=on
go install

# Tests, both unit tests and integration tests under /tests
echo "==> Running jk tests"
go test -v ./...
