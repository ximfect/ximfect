/* generic utilities */

package environ

import (
	"os"
)

// PrepareAppdata creates an APPDATA/ximfect folder and it's structure if it doesn't exist
func PrepareAppdata() {
	if _, err := os.Stat(Appdata); os.IsNotExist(err) {
		os.Mkdir(Appdata, os.ModePerm)
		os.Mkdir(AppdataPath("effects"), os.ModePerm)
	}
}

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
