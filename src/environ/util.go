/* generic utilities */

package environ

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
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

// SaveImage saves the given image
func SaveImage(filename string, img *image.RGBA) error {
	var format string
	if strings.Contains(filename, ".") {
		split := strings.Split(filename, ".")
		format = split[len(split)-1]
	} else {
		filename += ".png"
		format = "png"
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "png":
		err = png.Encode(file, img)
		if err != nil {
			return err
		}
	case "jpeg":
		err = jpeg.Encode(file, img, nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	return nil
}
