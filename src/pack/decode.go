package pack

import (
	// "bytes"
	// "compress/gzip"
	"errors"
)

const (
	filehdr = "XIMPKG"
)

type fileMap map[string][]byte

// Package describes a ximfect package (.xpk file)
type Package struct {
	Name string
	Files fileMap
}

// GetPackage reads and decodes a Package
func GetPackage(raw []byte) (*Package, error) {
	/*
	r, err := gzip.NewReader(bytes.NewBuffer(raw))
	if err != nil {
		return nil, err
	}
	src := make([]byte, 0xFFFF)
	n, err := r.Read(src)
	if err != nil {
		return nil, err
	}
	src = src[0:n]
	*/
	src := raw
	if string(src[0:6]) != filehdr {
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