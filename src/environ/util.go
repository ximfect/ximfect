/* generic utilities */

package environ

import (
	"os"
)

func EnsureDir(path string) {
	target, err := os.Stat(path)
	if err != nil || !target.IsDir() {
		_ = os.Mkdir(path, os.ModePerm)
	}
}

func EnsureAppdata() {
	EnsureDir(ProgramData)
	EnsureDir(DataPath("effects"))
	EnsureDir(DataPath("libs"))
	EnsureDir(DataPath("generators"))
	EnsureDir(DataPath("logs"))
}
