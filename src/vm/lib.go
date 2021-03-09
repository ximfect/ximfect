/* lib definitions */

package vm

import (
	lua "github.com/yuin/gopher-lua"
	"ximfect/environ"
)

// LibMetadata holds additional information about a Lib
type LibMetadata struct {
	Name    string
	Version string
	ID      string
	Author  string
	Desc    string
}

// Lib represents a loaded library
type Lib struct {
	Metadata *LibMetadata
	Files    []string
	Dir      string
}

// NewLib returns a pointer to a Lib
func NewLib(meta *LibMetadata, files []string, dir string) *Lib {
	return &Lib{meta, files, dir}
}

// Apply applies this Lib to the given VM
func (l *Lib) Apply(vm *lua.LState) error {
	var err error
	for _, file := range l.Files {
		err = vm.DoFile(environ.Combine(l.Dir, file))
		if err != nil {
			return err
		}
	}

	return nil
}