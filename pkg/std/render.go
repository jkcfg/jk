package std

import (
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
)

func (std *Std) render(pluginInfo string, params []byte) ([]byte, error) {
	renderer, err := std.plugins.GetRenderer(pluginInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching %s", pluginInfo)
	}

	result, err := renderer.Render(params)
	if err != nil {
		return nil, err
	}

	encoder := unicode.UTF16(NativeEndian, unicode.IgnoreBOM).NewEncoder()
	return encoder.Bytes(result)
}
