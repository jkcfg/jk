package resolve

import (
	"github.com/jkcfg/jk/pkg/vfs"
)

// MagicImporter handles importing "magic" modules, that is modules
// that are calculated wherever they are imported.
type MagicImporter struct {
	Specifier string
	Generate  func(vfs.Location) ([]byte, string)
	Public    bool // can this be imported from any module?
}

// Import implements the Importer interface for MagicImporter.
func (m *MagicImporter) Import(base vfs.Location, specifier, referrer string) ([]byte, vfs.Location, []Candidate) {
	if m.Specifier == specifier {
		if !m.Public && !IsStdModule(referrer) {
			return nil, vfs.Nowhere, nil
		}
		source, p := m.Generate(base)
		return source, vfs.Location{Vfs: vfs.Empty, Path: p}, nil
	}
	return nil, vfs.Nowhere, nil
}
