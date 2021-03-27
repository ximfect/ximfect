/* effect loading */

package vm

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"ximfect/environ"

	"gopkg.in/yaml.v2"
)

// LoadEffect loads an Effect from the given directory with the given id.
func LoadEffect(path, id string) (*Effect, error) {
	var err error

	dir := environ.Combine(path, id)
	meta := new(EffectMetadata)

	// Populate ID as it's not present in effect.yml
	meta.ID = id
	// Open meta file
	metaSource, err := os.Open(environ.Combine(dir, "effect.yml"))
	if err != nil {
		return nil, fmt.Errorf("error while loading metadata: %v", err)
	}
	// Create meta decoder & read meta
	metaDecoder := yaml.NewDecoder(metaSource)
	err = metaDecoder.Decode(meta)
	if err != nil {
		return nil, fmt.Errorf("error while reading metadata: %v", err)
	}

	// Return loaded effect
	return NewEffect(meta, dir), nil
}

// LoadAppdataEffect does what Load does, but path is always APPDATA
func LoadAppdataEffect(id string) (*Effect, error) {
	return LoadEffect(environ.DataPath("effects"), id)
}

// LoadLib loads a Lib from the given directory with the given id.
func LoadLib(path, id string) (*Lib, error) {
	dir := environ.Combine(path, id)
	metaPath := environ.Combine(dir, "lib.yml")

	var (
		err           error
		metaSource    *os.File
		metaDecoder   *yaml.Decoder
		meta          *LibMetadata = new(LibMetadata)
		filesAll      []os.FileInfo
		filesFiltered []string = []string{}
		fileName      string
	)

	// Populate ID as it's not present in effect.yml
	meta.ID = id
	// Open metadata file
	metaSource, err = os.Open(metaPath)
	if err != nil {
		return nil, fmt.Errorf("error loading metadata: %v", err)
	}
	// Create decoder & read meta
	metaDecoder = yaml.NewDecoder(metaSource)
	err = metaDecoder.Decode(meta)
	if err != nil {
		return nil, fmt.Errorf("error reding metadata: %v", err)
	}
	// Get list of all files in lib's folder
	filesAll, err = ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error discovering library files: %v",
			err)
	}
	// Filter files to only scripts
	for _, file := range filesAll {
		fileName = file.Name()
		if strings.HasSuffix(fileName, ".lua") {
			filesFiltered = append(filesFiltered, fileName)
		}
	}

	return NewLib(meta, filesFiltered, dir), nil
}

// LoadAppdataLib does what Load does, but path is always APPDATA
func LoadAppdataLib(id string) (*Lib, error) {
	return LoadLib(environ.DataPath("libs"), id)
}
