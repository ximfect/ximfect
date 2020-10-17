package pack

import "errors"

type fileMap map[string][]byte

// Package describes a ximfect package (.xpk file)
type Package struct {
	Name string
	Files fileMap
}

func GetPackage(src []byte) (*Package, error) {
	if !(src[0] == 'X' && src[1] == 'I' &&
		src[2] == 'M' && src[3] == 'P' &&
		src[4] == 'K' && src[5] == 'G') {
		return nil, errors.New("missing file header")
	}
	var (
		name    []byte
		hasName bool

		filename    []byte
		hasFilename bool
		contents    []byte
		files       fileMap = make(fileMap)
	)
	for _, c := range src[6:] {
		if !hasName {
			if c != 0x00 {
				name = append(name, c)
			} else {
				hasName = true
			}
		} else {
			if !hasFilename {
				if c != 0x00 {
					filename = append(filename, c)
				} else {
					hasFilename = true
				}
			} else {
				if c != 0x00 {
					contents = append(contents, c)
				} else {
					files[string(filename)] = contents
					hasFilename = false
					filename = []byte{}
					contents = []byte{}
				}
			}
		}
	}
	return &Package{string(name), files}, nil
}