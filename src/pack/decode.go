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
	Name  string
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
	// make sure the file header is there (prevents us from possibly reading
	// random files, most of the time at least)
	if string(src[0:6]) != filehdr {
		return nil, errors.New("missing file header")
	}
	// context variables
	var (
		ctx []byte

		packageName    string
		hasPackageName bool
		fileName       string
		hasFileName    bool
		files          = make(fileMap)
	)

	// for every byte in the source byte slice
	for _, b := range src[6:] {
		// if we don't have a package name
		if !hasPackageName {
			// if this byte is 0x00
			if b == 0x00 {
				// we have our package name
				hasPackageName = true
				packageName = string(ctx)
				// clear the context
				ctx = []byte{}
				// go to next byte
				continue
			}
			// if we have a package name, but not a filename
		} else if !hasFileName {
			// if this byte is 0x00
			if b == 0x00 {
				// we have our filename
				hasFileName = true
				fileName = string(ctx)
				// clear the context
				ctx = []byte{}
				// go to next byte
				continue
			}
			// if we have a package name and a filename
		} else {
			// if this byte is 0x00
			if b == 0x00 {
				// we have our file contents
				files[fileName] = ctx
				// reset the context
				hasFileName = false
				fileName = ""
				// go to next byte
				continue
			}
		}
		// append this byte to context
		ctx = append(ctx, b)
	}
	return &Package{packageName, files}, nil
}
