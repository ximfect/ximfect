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
