jk run -f %b/params.json -d %b.js | sed 's#^\(.*"path": "\).*\(/jk/tests/.*\)$#\1\2#'
