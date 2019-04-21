#!/bin/bash -e

std=$(dirname $0)
root=${std}/..

for f in `ls ${std}/*.fbs`; do
  flatc --go -o ${root}/pkg ${f}
done

flatc --ts --gen-all --no-ts-reexport -o ${std} ${std}/__std.fbs
flatc --js --gen-onefile --gen-all --es6-js-export -o ${std} ${std}/__std.fbs
