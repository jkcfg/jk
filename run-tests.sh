#!/bin/sh -e

export GO111MODULE=on

echo
echo "==> Running std tests"
(
  cd std
  npm test
)

echo
echo "==> Running eslint on std"
(
  cd std
  npm run lint
)

echo
echo "==> Running eslint on tests"
(
  cd std
  npx eslint -c .eslintrc ../tests/*.js
)

echo
echo "==> Running eslint on examples"
(
  cd std
  npx eslint -c .eslintrc ../examples
)

echo
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
echo
echo "==> Running jk tests"
go test -v ./...

echo
echo "==> Checking committed generated files are up to date"
git diff --exit-code

