package std

import (
	"fmt"
)

func MakeResourceModule(basePath string) ([]byte, string) {
	return []byte(fmt.Sprintf(`
import std from '@jkcfg/std';

const base = %q;

function resource(path, ...rest) {
  return std.read(base +'/' + path, ...rest);
}

export default resource;
`,
		basePath)), "resource:" + basePath
}
