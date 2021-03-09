package pack

import (
	"os"
	"ximfect/environ"
)

// UnpackTo reads the given Package and unpacks it to the given directory
func UnpackTo(pkg *Package, dest string) error {
	_ = os.Mkdir(dest, os.ModePerm)
	for name, content := range pkg.Files {
		file, err := os.Create(environ.Combine(dest, name))
		if err != nil {
			return err
		}
		_, err = file.Write(content)
		if err != nil {
			return err
		}
		_ = file.Close()
	}
	return nil
}