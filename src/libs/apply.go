/* loading libs into VMs */

package libs

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"ximfect/environ"
)

// ApplyLib applies the given Lib to the given VM
func ApplyLib(vm *lua.LState, lib *Lib) error {
	for _, file := range lib.Files {
		filePath := environ.Combine(lib.Dir, file)
		err := vm.DoFile(filePath)
		if err != nil {
			return fmt.Errorf("error while running library file `%s`: %v",
				file, err)
		}
	}
	return nil
}
