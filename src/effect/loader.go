package effect

import (
	"fmt"
	"strings"
	"ximfect/cfg"
	"ximfect/environ"
)

// Load loads an Effect from the given directory with the given id.
func Load(path, id string) (*Effect, error) {
	dir := environ.Combine(path, id)
	metaPath := environ.Combine(dir, "effect.meta")
	scriptPath := environ.Combine(dir, "effect.js")

	var (
		err        error
		metaSource string
		metaParsed cfg.Config
		script     string
		meta       *Metadata
		fx         *Effect
	)

	metaSource, err = environ.LoadTextfile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("error while loading metadata: %v", err)
	}
	script, err = environ.LoadTextfile(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("error while loading script: %v", err)
	}

	metaParsed = cfg.Parse(metaSource)

	var (
		name       string
		version    string
		author     string
		desc       string
		preloadRaw string
		preload    []string
		ok         bool
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
	preloadRaw, ok = metaParsed["preload"]
	if ok {
		preload = strings.Split(preloadRaw, " ")[1:]
	} else {
		preload = []string{}
	}

	meta = &(Metadata{name, version, id, author, desc, preload})
	fx = NewEffect(meta, script)
	return fx, nil
}

// LoadFromAppdata does what Load does, but path is always APPDATA
func LoadFromAppdata(id string) (*Effect, error) {
	return Load(environ.AppdataPath("effects"), id)
}
