// +build dev

package lib

import "net/http"

// Assets is the content of the std/build directory.
var Assets http.FileSystem = http.Dir("../../../std/dist")
