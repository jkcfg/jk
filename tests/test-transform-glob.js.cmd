jk transform --stdout -c '({ number }) => ({ plusone: number + 1 })' ./test-transform-files/*.yaml
# This tests that
#  1. all files are processed
#  2. the transformed documents from multidoc files are combined with
#     documents into one multidoc output
