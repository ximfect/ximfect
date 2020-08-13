/* lib definitions */

package libs

// Metadata holds additional information about a Lib
type Metadata struct {
	Name    string
	Version string
	ID      string
	Author  string
	Desc    string
}

// Lib represents a loaded library
type Lib struct {
	Metadata *Metadata
	Files    []string
	Dir      string
}

// NewLib returns a pointer to a Lib
func NewLib(meta *Metadata, files []string, dir string) *Lib {
	return &Lib{meta, files, dir}
}
