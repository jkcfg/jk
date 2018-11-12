#!/bin/sh -e

export GO111MODULE=on

echo "==> Running std tests"
(
  cd std
  npm test
)

echo "==> Running go linters"
gometalinter --tests --vendor --disable-all --deadline=600s \
    --enable=misspell \
    --enable=vet \
    --enable=ineffassign \
    --enable=gofmt \
    --enable=deadcode \
    --enable=golint \
    ./...

# Tests, both unit tests and integration tests under /tests
echo "==> Running jk tests"
go test -v ./...

echo "==> Checking committed generated files are up to date"
git diff --exit-code

