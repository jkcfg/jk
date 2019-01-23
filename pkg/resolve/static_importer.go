package resolve

const (
	staticRule = "built-in"
)

// StaticImporter is an importer mapping an import specifier to a static string.
type StaticImporter struct {
	Specifier string
	Source    []byte
}

// Import implements importer.
func (si *StaticImporter) Import(basePath, specifier, referrer string) ([]byte, string, []Candidate) {
	candidate := []Candidate{{specifier, staticRule}}
	if si.Specifier == specifier {
		return si.Source, specifier, candidate
	}
	return nil, "", candidate
}
