/* generic utilities */

package environ

import (
	"os"
)

// LoadTextfile opens and reads a text file, returns it's contents as a string
func LoadTextfile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buffer := make([]byte, 0xFFFF)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	out := ""
	for _, b := range buffer {
		if b == 0 {
			break
		}
		out += string(b)
	}

	return out, nil
}

// LoadRawfile opens and reads a text file, returns it's contents as a string
func LoadRawfile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()

	buffer := make([]byte, 0xFFFF)
	amt, err := file.Read(buffer)
	if err != nil {
		return []byte{}, err
	}

	return buffer[0:amt], nil
}

func EnsureAppdata() {
	var err error

	ensureDir := func(path string) {
		_, err = os.Stat(path)
		if err != nil {
			os.Mkdir(path, os.ModePerm)
		}
	}

	ensureDir(Appdata)
	ensureDir(AppdataPath("effects"))
	ensureDir(AppdataPath("libs"))
	ensureDir(AppdataPath("logs"))
}