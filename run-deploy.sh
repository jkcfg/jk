#!/bin/bash

set -e
set -x

# Only deploy from the Linux CI
os=`go env GOOS`
if [ $os != "linux" ]; then
  exit 0
fi

function website_upload() {
  src=$1
  dst=$2

  if [ -z "$CI" ]; then
    # Not on CI
    git_url="git@github.com:jkcfg/jkcfg.github.io.git"
  else
    # CI
    git_url="https://${GITHUB_TOKEN}@github.com/jkcfg/jkcfg.github.io.git"

    echo "==> setting up git"
    git config --global user.email "damien.lespiau+jkbot@gmail.com"
    git config --global user.name "jkbot"
  fi

  echo "==> cloning site"
  git clone $git_url deploy

  echo "==> deploying $src to $dst"
  rm -rf deploy/$dst
  mkdir -p $(dirname deploy/$dst)
  cp -r $src deploy/$dst
  (cd deploy && git add . && git commit -m "automated publication" && git push)
  rm -rf deploy
}

website_upload std/docs reference/std/latest
