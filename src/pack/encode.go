package pack

import (
	"io/ioutil"
	"path/filepath"
	// "bytes"
	// "compress/gzip"
	"ximfect/environ"
)

// PackageFile represents a file in a Package
type PackageFile struct {
	Name     string
	Path     string
	Contents []byte
}

// GetPacked creates a Package with the given PackageFiles
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
	/*
	comp := make([]byte, 0xFFFF)
	w := gzip.NewWriter(bytes.NewBuffer(comp))
	n, err := w.Write(out)
	if err != nil {
		return []byte{}, err
	}
	return comp[0:n], nil
	*/
	return out, nil
}

// GetPackedDirectory returns the given folder as a Package
func GetPackedDirectory(path string) ([]byte, error) {
	files, err := GetPackageDirectory(path)
	if err != nil {
		return nil, err
	}
	return GetPacked(filepath.Base(path), files)
}

// GetPackageDirectory returns a directory's contents as PackageFiles
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