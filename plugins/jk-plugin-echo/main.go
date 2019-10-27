package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/jkcfg/jk/pkg/plugin/renderer"
)

// Echo outputs its input.
type Echo struct {
	log hclog.Logger
}

// Render implements renderer.Renderer.
func (h *Echo) Render(input []byte) ([]byte, error) {
	h.log.Debug("debug message from echo plugin")
	return input, nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Info,
		Output: os.Stderr,
	})

	r := &Echo{
		log: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"renderer": &renderer.Plugin{Impl: r},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: renderer.RendererV1,
		Plugins:         pluginMap,
	})
}
