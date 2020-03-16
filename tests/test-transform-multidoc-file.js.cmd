rm -rf ./%d/test-transform-files
jk transform -c '({ number }) => ({ plusone: number + 1 })' ./test-transform-files/numbers.yaml -o %d
