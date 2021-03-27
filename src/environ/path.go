/* path utilities */

package environ

import "strings"

// Combine combines a base path and a slice of folders/files to add
func Combine(base string, path ...string) string {
	combined := append([]string{base}, path...)
	return strings.Join(combined, PathSep)
}

// DataPath returns a path in the ProgramData directory
func DataPath(path ...string) string {
	return Combine(ProgramData, path...)
}
