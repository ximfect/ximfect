package effect

import (
	"fmt"
	"os"
	"ximfect/cfg"
	"ximfect/environ"
)

// Load loads an Effect from the given directory with the given id.
func Load(path, id string) (*Effect, error) {
	dir := environ.Combine(path, id)
	metaPath := environ.Combine(dir, "effect.meta")
	scriptPath := environ.Combine(dir, "effect.js")

	var (
		err          error
		metaFile     *os.File
		scriptFile   *os.File
		metaBuffer   []byte
		scriptBuffer []byte
		metaSource   string
		metaParsed   cfg.Config
		script       string
		meta         *Metadata
		fx           *Effect
	)

	metaFile, err = os.Open(metaPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening metadata file: %v", err)
	}
	scriptFile, err = os.Open(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening script file: %v", err)
	}

	defer metaFile.Close()
	defer scriptFile.Close()

	metaBuffer = make([]byte, 0xFFFF)
	_, err = metaFile.Read(metaBuffer)
	if err != nil {
		return nil, fmt.Errorf("error while reading metadata: %v", err)
	}
	metaSource = ""
	for _, ch := range metaBuffer {
		if ch == 0 {
			break
		}
		metaSource += string(ch)
	}

	scriptBuffer = make([]byte, 0xFFFF)
	_, err = scriptFile.Read(scriptBuffer)
	if err != nil {
		return nil, fmt.Errorf("error while reading script: %v", err)
	}
	script = ""
	for _, ch := range scriptBuffer {
		if ch == 0 {
			break
		}
		script += string(ch)
	}

	metaParsed = cfg.Parse(metaSource)

	var (
		name    string
		version string
		author  string
		desc    string
		ok      bool
	)

	name, ok = metaParsed["name"]
	if !ok {
		return nil, fmt.Errorf(
			"error while applying metadata: could not find required field `name`")
	}
	version, ok = metaParsed["version"]
	if !ok {
		return nil, fmt.Errorf(
			"error while applying metadata: could not find required field `version`")
	}
	author, ok = metaParsed["author"]
	if !ok {
		return nil, fmt.Errorf(
			"error while applying metadata: could not find required field `author`")
	}
	desc, ok = metaParsed["desc"]
	if !ok {
		return nil, fmt.Errorf(
			"error while applying metadata: could not find required field `desc`")
	}

	meta = &(Metadata{name, version, id, author, desc})
	fx = NewEffect(meta, script)
	return fx, nil
}

// LoadFromAppdata does what Load does, but path is always APPDATA
func LoadFromAppdata(id string) (*Effect, error) {
	return Load(environ.AppdataPath("effects"), id)
}
