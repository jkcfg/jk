package std

import (
	"golang.org/x/text/encoding/unicode"
)

// NativeEndian is the platform-specific endianness. We need this
// because certain things in JS, e.g., Uint16Array, use "native
// endianness".
const NativeEndian = unicode.LittleEndian
