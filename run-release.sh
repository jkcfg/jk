#!/bin/bash

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

binary=jk-`go env GOOS`-`go env GOARCH`
mv jk $binary

echo "==> Uploading $binary"
upload $binary

echo "==> Uploading $binary.sha256"
shasum -a 256 $binary > $binary.sha256
upload $binary.sha256
