#!/bin/bash

tag=$1
user=dlespiau
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

echo "==> Uploading jk-linux-amd64"
run github-release upload \
    --user $user \
    --repo $repo \
    --tag $tag \
    --name "jk-linux-amd64" \
    --file jk
