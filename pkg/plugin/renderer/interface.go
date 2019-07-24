package renderer

import "github.com/hashicorp/go-plugin"

// RendererV1 is the first version of the render plugin interface.
var RendererV1 = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "JK_PLUGIN",
	MagicCookieValue: "renderer",
}

// Renderer is a plugin that outputs one or more structured objects.
//
// Maps are sent and received as json objects in byte arrays. It's a lot simpler
// to send a byte array through RPC than trying to have the message encoding try
// to serialize a map, especially because we can hand over the byte array to v8
// directory and unmarshal the JSON there.
//
// Both the input and the result byte arrays are JSON documents serialized as
// utf8 strings.
type Renderer interface {
	Render(input []byte) ([]byte, error)
}
