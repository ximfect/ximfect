/* generic utilities */

package environ

import (
	"os"
)

func EnsureDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return err
}

func EnsureAppdata() error {
	err := EnsureDir(ProgramData)
	if err != nil {
		return err
	}
	err = EnsureDir(DataPath("effects"))
	if err != nil {
		return err
	}
	err = EnsureDir(DataPath("libs"))
	if err != nil {
		return err
	}
	err = EnsureDir(DataPath("generators"))
	if err != nil {
		return err
	}
	err = EnsureDir(DataPath("logs"))
	return err
}
