package plugin

import (
	"runtime"
)

// Info holds plugin metadata.
type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// Binaries maps `go env GOOS`-`go env GOARCH` strings to binary names.
	// eg. "linux-amd64" -> "https://jkcfg.github.io/plugins/render/echo/0.1.0/jk-render-echo-linux-amd64"
	Binaries map[string]string `json:"binaries"`
}

func (i *Info) binary() string {
	k := runtime.GOOS + "-" + runtime.GOARCH
	return i.Binaries[k]
}
