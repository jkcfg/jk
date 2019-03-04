package resolve

// MagicImporter handles importing "magic" modules, that is modules
// that are calculated wherever they are imported.
type MagicImporter struct {
	Specifier string
	Generate  func(string) ([]byte, string)
}

// Import implements the Importer interface for MagicImporter.
func (m *MagicImporter) Import(basePath, specifier, referrer string) ([]byte, string, []Candidate) {
	if m.Specifier == specifier {
		source, path := m.Generate(basePath)
		return source, path, nil
	}
	return nil, "", nil
}
