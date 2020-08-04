package environ

import (
	"os"
)

// Sep is a constant string path separator
const Sep string = string(os.PathSeparator)

// Appdata is a constant path of APPDATA
var Appdata string = os.Getenv("APPDATA") + Sep + "ximfect"

// Combine combines a base path and a slice of folders/files to add
func Combine(base string, path ...string) string {
	out := base
	for _, name := range path {
		out += Sep + name
	}
	return out
}

// AppdataPath returns a path in APPDATA
func AppdataPath(path ...string) string {
	return Combine(Appdata, path...)
}
