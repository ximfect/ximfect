/* generic utilities */

package environ

import (
	"os"
)

// LoadTextfile opens and reads a text file, returns it's contents as a string
func LoadTextfile(path string) (string, error) {
	out := []byte{}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buffer := make([]byte, 0xFFFF)
	n, err := file.Read(buffer)
	total := n
	if err != nil {
		return "", err
	}
	out = append(out, buffer...)

	for n == 0xFFFF {
		out = append(out, buffer...)
		buffer = make([]byte, 0xFFFF)
		_, err = file.Seek(int64(n), 1)
		if err != nil {
			return "", err
		}
		n, err = file.Read(buffer)
		if err != nil {
			return "", err
		}
		total += n
	}

	return string(out[0:total]), nil
}

// LoadRawfile opens and reads a text file, returns it's contents as a []byte
func LoadRawfile(path string) ([]byte, error) {
	out := []byte{}
	file, err := os.Open(path)
	if err != nil {
		return out, err
	}

	defer file.Close()

	buffer := make([]byte, 0xFFFF)
	n, err := file.Read(buffer)
	total := n
	if err != nil {
		return out, err
	}
	out = append(out, buffer...)

	for n == 0xFFFF {
		out = append(out, buffer...)
		buffer = make([]byte, 0xFFFF)
		_, err = file.Seek(int64(n), 1)
		if err != nil {
			return out, err
		}
		n, err = file.Read(buffer)
		if err != nil {
			return out, err
		}
		total += n
	}

	return out[0:total], nil
}

func EnsureAppdata() {
	var err error

	ensureDir := func(path string) {
		_, err = os.Stat(path)
		if err != nil {
			_ = os.Mkdir(path, os.ModePerm)
		}
	}

	ensureDir(Appdata)
	ensureDir(AppdataPath("effects"))
	ensureDir(AppdataPath("libs"))
	ensureDir(AppdataPath("logs"))
}
