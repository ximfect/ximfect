package libs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"ximfect/cfg"
	"ximfect/environ"
)

// Load loads a Lib from the given directory with the given id.
func Load(path, id string) (*Lib, error) {
	dir := environ.Combine(path, id)
	metaPath := environ.Combine(dir, "lib.meta")

	var (
		err           error
		metaSource    string
		metaParsed    cfg.Config
		meta          *Metadata
		filesAll      []os.FileInfo
		filesFiltered []string = []string{}
		fileName      string
	)

	metaSource, err = environ.LoadTextfile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("error while loading metadata: %v", err)
	}
	metaParsed = cfg.Parse(metaSource)

	filesAll, err = ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error while discovering library files: %v",
			err)
	}
	for _, file := range filesAll {
		fileName = file.Name()
		if strings.HasSuffix(fileName, ".js") {
			filesFiltered = append(filesFiltered, fileName)
		}
	}

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
	return NewLib(meta, filesFiltered, dir), nil
}

// LoadFromAppdata does what Load does, but path is always APPDATA
func LoadFromAppdata(id string) (*Lib, error) {
	return Load(environ.AppdataPath("libs"), id)
}
