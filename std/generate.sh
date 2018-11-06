#!/bin/bash -e

std=$(dirname $0)
root=${std}/..

for f in `ls ${std}/*.fbs`; do
  flatc --js --es6-js-export -o ${std} ${f}
  flatc --go -o ${root}/pkg ${f}
done
