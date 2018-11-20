#!/bin/sh -e

export GO111MODULE=on

echo "==> Running std tests"
(
  cd std
  npm test
)

echo "==> Running eslint on tests"
(
  cd std
  npx eslint -c .eslintrc ../tests/*.js
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

# We allow pkg/__std/lib/assets_vfsdata.go to be re-generated, but with only
# timestamp changed.
vfsdata=pkg/__std/lib/assets_vfsdata.go
changed=`git diff $vfsdata | grep ^[-+] | grep -v [ab]/$vfsdata | grep -v modTime: | wc -l`
[ $changed -eq 0 ] && git checkout $vfsdata

git diff --exit-code

