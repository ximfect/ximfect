package libs

import (
	"fmt"
	"ximfect/environ"

	"github.com/robertkrimen/otto"
)

// ApplyLib applies the given Lib to the given VM
func ApplyLib(vm *otto.Otto, lib *Lib) error {
	for _, file := range lib.Files {
		filePath := environ.Combine(lib.Dir, file)
		libFile, err := environ.LoadTextfile(filePath)
		if err != nil {
			return fmt.Errorf("error while loading library file `%s`: %v",
				file, err)
		}
		_, err = vm.Run(libFile)
		if err != nil {
			return fmt.Errorf("error while running library file `%s`: %v",
				file, err)
		}
	}
	return nil
}
