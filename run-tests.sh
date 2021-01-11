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
golangci-lint run --no-config --timeout=600s \
  --disable-all --enable=deadcode --enable=golint --enable=varcheck \
  --enable=structcheck --enable=dupl --enable=ineffassign \
  --enable=interfacer --enable=govet --enable=gofmt --enable=misspell \
  ./...

# these are unused linters at present; some may fail.
# --enable=unconvert --enable=megacheck --enable=errcheck --enable=maligned
# --enable=goconst --enable=gosec

# Tests, both unit tests and integration tests under /tests
echo
echo "==> Running jk tests"
go test -v ./...

if [ -n "$CI" ]; then
  echo
  echo "==> Checking committed generated files are up to date"
  git diff --exit-code
fi
