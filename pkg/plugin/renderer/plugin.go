package renderer

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

// RPCResponse is the Render response on the wire.
type RPCResponse struct {
	Data []byte
	Err  error
}

// RPCClient is a renderer implemented with golang's built-in RPC mechanism.
type RPCClient struct{ client *rpc.Client }

// Render implements the Renderer interface.
func (r *RPCClient) Render(input []byte) ([]byte, error) {
	var resp RPCResponse
	err := r.client.Call("Plugin.Render", input, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "render")
	}

	return resp.Data, resp.Err
}

// RPCServer is the RPC server that RPCClient talks to, conforming to the
// requirements of net/rpc.
type RPCServer struct {
	Impl Renderer
}

// Render is the net/rpc implementation of the Renderer interface.
func (s *RPCServer) Render(input []byte, resp *RPCResponse) error {
	data, err := s.Impl.Render(input)
	resp.Data = data
	resp.Err = err
	return nil
}

// Plugin is the implementation of plugin.Plugin so we can serve/consume this
// interface.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return RPCClient for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type Plugin struct {
	Impl Renderer
}

// Server implements the plugin.Plugin interface.
func (p *Plugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

// Client implements the plugin.Plugin interface.
func (Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}
