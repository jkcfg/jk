#!/bin/sh -e

echo "==> Running std tests"
(
  cd std
  npm test
)

# Tests, both unit tests and integration tests under /tests
echo "==> Running jk tests"
GO111MODULE=on go test -v ./...

echo "==> Checking committed generated files are up to date"
git diff --exit-code

