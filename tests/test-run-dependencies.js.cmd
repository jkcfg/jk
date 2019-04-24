jk run -f %t/params.json -d %t.js | sed 's#^\(.*"path": "\).*\(/jk/tests/.*\)$#\1\2#'
