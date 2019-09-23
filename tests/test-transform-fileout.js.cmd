rm -rf ./test-transform-fileout.got
jk transform -o ./test-transform-fileout.got -c 'v => v + 1' ./test-transform-files/*.json
# NB:
# * --stdout=false, so expected to write a file
# * -o won't be set automatically for a .cmd file
