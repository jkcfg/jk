rm -rf ./test-transform-good-overwrite.got
mkdir -p ./test-transform-good-overwrite.got/test-transform-files
echo '1' > ./test-transform-good-overwrite.got/test-transform-files/one.json
jk transform --overwrite -o ./test-transform-good-overwrite.got/ -c 'v => v + 1' ./test-transform-files/one.json
