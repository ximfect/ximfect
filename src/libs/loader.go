/* lib loading */

package libs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"ximfect/environ"

	"gopkg.in/yaml.v2"
)

// Load loads a Lib from the given directory with the given id.
func Load(path, id string) (*Lib, error) {
	dir := environ.Combine(path, id)
	metaPath := environ.Combine(dir, "lib.yml")

	var (
		err           error
		metaSource    *os.File
		metaDecoder   *yaml.Decoder
		meta          *Metadata = new(Metadata)
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

// LoadFromAppdata does what Load does, but path is always APPDATA
func LoadFromAppdata(id string) (*Lib, error) {
	return Load(environ.AppdataPath("libs"), id)
}
