#!/bin/bash

set -e

tag=$1
user=jkcfg
repo=jk
pkg=github.com/$user/$repo

function run() {
    if command -v github-release; then
        github-release "$@"
    else
        docker run -e GITHUB_TOKEN -v "$(pwd)":/go/src/$pkg quay.io/justkidding/build github-release "$@"
    fi
}

echo "==> Checking package.json is up to date"
version=$(./jk run std/version.jk)
if [ "$version" != "$tag" ]; then
    echo "error: releasing $tag but std/package.json references $version"
    exit 1
fi

echo "==> Creating $tag release"
run release \
    --user $user \
    --repo $repo \
    --tag $tag

function upload() {
    file=$1
    run upload \
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

echo "==> Uploading npm module"
if [ -z "$NPM_TOKEN" ]; then
    echo "error: NPM_TOKEN needs to be defined for  pushing npm modules"
    exit 1
fi
echo '//registry.npmjs.org/:_authToken=${NPM_TOKEN}' > @jkcfg/std/.npmrc
(cd @jkcfg/std && npm publish)
