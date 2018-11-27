#!/bin/bash

tag=$1
user=justkidding-config
repo=jk
pkg=github.com/$user/$repo

function run() {
    docker run -e GITHUB_TOKEN -v "$(pwd)":/go/src/$pkg quay.io/justkidding/build "$@"
}

echo "==> Creating $tag release"
run github-release release \
    --user $user \
    --repo $repo \
    --tag $tag

function upload() {
    file=$1
    run github-release upload \
        --user $user \
        --repo $repo \
        --tag $tag \
        --name $file \
        --file $file

}

binary=jk-linux-amd64
mv jk $binary

echo "==> Uploading $binary"
upload $binary

echo "==> Uploading $binary.sha256"
shasum -a 256 $binary > $binary.sha256
upload $binary.sha256
