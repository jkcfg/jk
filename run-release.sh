#!/bin/bash

tag=$1
user=dlespiau
repo=jk
pkg=github.com/$user/$repo

function run() {
    docker run -e GITHUB_TOKEN -v "$(pwd)":/go/src/$pkg quay.io/justkidding/build "$@"
}

run github-release release \
    --user $user \
    --repo $repo \
    --tag $tag \
    --name $tag

run github-$ upload \
    --user $user \
    --repo $repo \
    --tag $tag \
    --name "jk-linux-amd64" \
    --file jk
