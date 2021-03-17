package pack

// import (
// 	"bytes"
// 	"compress/gzip"
// )

// PackageFile represents a file in a Package
type PackageFile struct {
	Name     string
	Path     string
	Contents []byte
}

// GetPackedDirectory creates a named Package from the given directory
func GetPackedDirectory(name, path string) ([]byte, error) {
	// get files
	files, err := GetDirectoryFiles(path)
	if err != nil {
		return nil, err
	}
	// file header
	out := []byte{'X', 'I', 'M', 'P', 'K', 'G'}
	// package name
	out = append(out, []byte(name)...)
	out = append(out, 0x00)
	// package files
	for _, file := range files {
		// file name
		out = append(out, []byte(file.Name)...)
		out = append(out, 0x00)
		// file contents
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
