package effect

import (
	"fmt"
	"os"
	"ximfect/environ"

	"gopkg.in/yaml.v2"
)

// Load loads an Effect from the given directory with the given id.
func Load(path, id string) (*Effect, error) {
	dir := environ.Combine(path, id)
	metaPath := environ.Combine(dir, "effect.yml")
	scriptPath := environ.Combine(dir, "effect.js")

	var (
		err         error
		metaSource  *os.File
		metaDecoder *yaml.Decoder
		script      string
		meta        *Metadata
	)

	// Create empty *Metadata
	meta = new(Metadata)
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

// LoadFromAppdata does what Load does, but path is always APPDATA
func LoadFromAppdata(id string) (*Effect, error) {
	return Load(environ.AppdataPath("effects"), id)
}
