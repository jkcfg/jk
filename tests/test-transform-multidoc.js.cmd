jk transform --stdout -c '({ number }) => ({ plusone: number + 1 })' ./test-transform-files/numbers.yaml
# This confirms that a multidoc will get all documents read
