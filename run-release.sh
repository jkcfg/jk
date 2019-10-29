#!/bin/bash

set -e
set -x

tag=$1
user=jkcfg
repo=jk
pkg=github.com/$user/$repo

function docker_run() {
    docker run -e GITHUB_TOKEN -e NPM_TOKEN -v "$(pwd)":/go/src/$pkg quay.io/justkidding/build "$@"
}

function run() {
    cmd=$1
    if [[ $cmd != ./* ]] && command -v $cmd; then
        "$@"
    else
        docker_run $@
    fi
}

echo "==> Creating $tag release"
run github-release release \
    --user $user \
    --repo $repo \
    --tag $tag \
    || true

function upload() {
    file=$1
    run github-release upload \
        --user $user \
        --repo $repo \
        --tag $tag \
        --name $file \
        --file $file

}

os=`go env GOOS`
binary=jk-$os-`go env GOARCH`
mv jk $binary

echo "==> Uploading $binary"
upload $binary

echo "==> Uploading $binary.sha256"
shasum -a 256 $binary > $binary.sha256
upload $binary.sha256

# We can only upload the npm module once. Do it from the Linux CI.
if [ $os != "linux" ]; then
  exit 0
fi

echo "==> Checking package.json is up to date"
version=$(run ./$binary run std/version.jk)
if [ "$version" != "$tag" ]; then
    echo "error: releasing $tag but std/package.json references $version"
    exit 1
fi

echo "==> Uploading npm module"
if [ -z "$NPM_TOKEN" ]; then
    echo "error: NPM_TOKEN needs to be defined for  pushing npm modules"
    exit 1
fi
echo '//registry.npmjs.org/:_authToken=${NPM_TOKEN}' > @jkcfg/std/.npmrc
docker_run bash -c '(cd @jkcfg/std && npm publish)'
