package pack

import (
	"io/ioutil"
	"os"
	"ximfect/environ"
)

// UnpackTo reads the given Package and unpacks it to the given directory
func UnpackTo(pkg *Package, dest string) error {
	// make sure destination directory exists
	_ = os.Mkdir(dest, os.ModePerm)
	// for every file
	for name, content := range pkg.Files {
		// create the file
		file, err := os.Create(environ.Combine(dest, name))
		if err != nil {
			return err
		}
		// write the contents
		_, err = file.Write(content)
		if err != nil {
			return err
		}
		// close it
		_ = file.Close()
	}
	return nil
}

// GetPackageDirectory returns a directory's contents as PackageFiles
func GetDirectoryFiles(path string) ([]PackageFile, error) {
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, nil
	}
	files := []PackageFile{}
	var (
		name     string
		filepath string
		contents []byte
	)
	for _, entry := range list {
		if entry.IsDir() {
			continue
		}
		name = entry.Name()
		filepath = environ.Combine(path, name)
		contents, err = environ.LoadRawfile(filepath)
		if err != nil {
			return nil, nil
		}
		files = append(files, PackageFile{name, filepath, contents})
	}
	return files, nil
}
