/* effect loading */

package vm

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"ximfect/environ"

	"gopkg.in/yaml.v2"
)

// LoadEffect loads an Effect from the given directory with the given id.
func LoadEffect(path, id string) (*Effect, error) {
	dir := environ.Combine(path, id)
	metaPath := environ.Combine(dir, "effect.yml")
	scriptPath := environ.Combine(dir, "effect.lua")

	var (
		err         error
		metaSource  *os.File
		metaDecoder *yaml.Decoder
		meta        *EffectMetadata = new(EffectMetadata)
		script      string
	)

	// Populate ID as it's not present in effect.yml
	meta.ID = id
	// Open meta file
	metaSource, err = os.Open(metaPath)
	if err != nil {
		return nil, fmt.Errorf("error while loading metadata: %v", err)
	}
	// Load script file
	script, err = environ.LoadTextfile(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("error while loading script: %v", err)
	}
	// Create meta decoder & read meta
	metaDecoder = yaml.NewDecoder(metaSource)
	err = metaDecoder.Decode(meta)
	if err != nil {
		return nil, fmt.Errorf("error while reading metadata: %v", err)
	}

	// Return loaded effect
	return NewEffect(meta, script), nil
}

// LoadAppdataEffect does what Load does, but path is always APPDATA
func LoadAppdataEffect(id string) (*Effect, error) {
	return LoadEffect(environ.AppdataPath("effects"), id)
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
		return nil, fmt.Errorf("error while loading metadata: %v", err)
	}
	// Create decoder & read meta
	metaDecoder = yaml.NewDecoder(metaSource)
	err = metaDecoder.Decode(meta)
	if err != nil {
		return nil, fmt.Errorf("erro while reding metadata: %v", err)
	}
	// Get list of all files in lib's folder
	filesAll, err = ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error while discovering library files: %v",
			err)
	}
	// Filter files to only javascript
	for _, file := range filesAll {
		fileName = file.Name()
		if strings.HasSuffix(fileName, ".js") {
			filesFiltered = append(filesFiltered, fileName)
		}
	}

	return NewLib(meta, filesFiltered, dir), nil
}

// LoadAppdataLib does what Load does, but path is always APPDATA
func LoadAppdataLib(id string) (*Lib, error) {
	return LoadLib(environ.AppdataPath("libs"), id)
}
