package tool

import (
	"os"
	"strings"
)

// Argument represents a CLI argument.
type Argument struct {
	HasValue bool
	Value    string
}

// ArgMap is a string->Argument map.
type ArgMap map[string]Argument

// GetArgv reads os.Args and turns it into an ArgMap.
func GetArgv() ArgMap {
	out := make(ArgMap)
	var last string
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			last = arg[2:]
			out[last] = Argument{false, ""}
		} else {
			if !out[last].HasValue {
				out[last] = Argument{true, arg}
			}
		}
	}
	return out
}
