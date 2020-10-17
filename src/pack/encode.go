package pack

import (
	"io/ioutil"
	"path/filepath"
	"ximfect/environ"
)

type PackageFile struct {
	Name     string
	Path     string
	Contents []byte
}

func GetPacked(name string, files []PackageFile) ([]byte, error) {
	out := []byte{'X', 'I', 'M', 'P', 'K', 'G'}
	out = append(out, []byte(name)...)
	out = append(out, 0x00)
	for _, file := range files {
		out = append(out, []byte(file.Name)...)
		out = append(out, 0x00)
		out = append(out, file.Contents...)
		out = append(out, 0x00)
	}
	return out, nil
}

func GetPackedDirectory(path string) ([]byte, error) {
	files, err := GetPackageDirectory(path)
	if err != nil {
		return nil, err
	}
	return GetPacked(filepath.Base(path), files)
}

func GetPackageDirectory(path string) ([]PackageFile, error) {
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