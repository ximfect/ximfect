package effect

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"ximfect/environ"
)

// FileMap is a utility map of []byte -> string
type FileMap map[string][]byte

// PackedEffect represents an effect package
type PackedEffect struct {
	Name   string
	Header []byte
	Files  FileMap
}

// GetPacked reads the contents of an effect package and returns a PackedEffect
func GetPacked(src []byte) (*PackedEffect, error) {
	if !(src[0] == 'X' && src[1] == 'F' && src[2] == 'P') {
		return nil, errors.New("missing file header")
	}
	var (
		name    string
		last    []byte
		nameEnd int
	)
	for i, c := range src[3:] {
		if c != '\xFF' {
			last = append(last, c)
		} else {
			name = string(last)
			nameEnd = i+4
			break
		}
	}
	//fmt.Println("name", name)
	//fmt.Println("nameEnd", nameEnd)
	var (
		headerSrc     []byte
		filesSrc      []byte
	)
	for i, c := range src[nameEnd:] {
		if c != '\xFF' {
			headerSrc = append(headerSrc, c)
		} else {
			filesSrc = src[i+nameEnd+1:]
			break
		}
	}
	//fmt.Println(filesSrc)
	//fmt.Println(string(filesSrc))
	var (
		hasName  bool
		files    FileMap = make(FileMap)
		fileName string
		ctx      []byte
	)
	for _, c := range filesSrc {
		if !hasName {
			//fmt.Println("no name")
			if c != '\xFF' {
				ctx = append(ctx, c)
			} else {
				//fmt.Println("end")
				hasName = true
				fileName = string(ctx)
				ctx = []byte{}
			}
		} else {
			//fmt.Println("has name")
			if c != 0xFF {
				ctx = append(ctx, c)
			} else {
				//fmt.Println("end")
				files[fileName] = ctx
				ctx = []byte{}
				hasName = false
			}
		}
	}
	if hasName && len(ctx) > 0 {
		files[fileName] = ctx
	}
	return &PackedEffect{name, headerSrc, files}, nil
}

// Unpack fully unpacks an effect package
func Unpack(name string) error {
	src, err := environ.LoadRawfile(name)
	if err != nil {
		return err
	}

	packed, err := GetPacked(src)
	if err != nil {
		return err
	}

	basePath := environ.AppdataPath("effects", packed.Name)
	os.Mkdir(basePath, os.ModePerm)
	metadataPath := environ.Combine(basePath, "effect.yml")
	metadataFile, err := os.Create(metadataPath)
	if err != nil {
		return err
	}

	_, err = metadataFile.Write(packed.Header)
	if err != nil {
		return err
	}

	for fileName, fileContent := range packed.Files {
		file, err := os.Create(environ.Combine(basePath, fileName))
		if err != nil {
			return err
		}
		_, err = file.Write(fileContent)
		if err != nil {
			return err
		}
	}

	return nil
}

// Pack fully packs an effect into a package
func Pack(effect *Effect) ([]byte, error) {
	out := []byte{'X', 'F', 'P'}
	id := effect.Metadata.ID
	out = append(out, []byte(id)...)
	out = append(out, '\xFF')
	headerSrc, err := yaml.Marshal(effect.Metadata)
	if err != nil {
		return []byte{}, err
	}
	out = append(out, headerSrc...)
	out = append(out, '\xFF')
	basePath := environ.AppdataPath("effects", id)
	fileList := []string{"effect.js"}
	fileList = append(fileList, effect.Metadata.Preload...)
	for _, fileName := range fileList {
		out = append(out, []byte(fileName)...)
		out = append(out, '\xFF')
		file, err := environ.LoadRawfile(
			environ.Combine(
				basePath, fileName))
		if err != nil {
			return []byte{}, nil
		}
		out = append(out, file...)
		out = append(out, '\xFF')
	}
	return out, nil
}