rm -rf ./test-transform-parent-path.got/
mkdir -p ./test-transform-parent-path.got
echo '1' > ./test-transform-parent-path.got/one.yaml
cd testfiles && jk transform --overwrite -c 'v => v+1' ../test-transform-parent-path.got/one.yaml
