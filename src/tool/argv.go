/* os.Args parsing */

package tool

import (
	"strings"
)

// Argument represents a CLI argument.
type Argument struct {
	IsValue   bool
	Value     string
	BoolValue bool
}

// ArgMap is a string->Argument map.
type ArgMap map[string]Argument

// ArgSlice is an Argument slice.
type ArgSlice []Argument

// ArgumentList turns a slice of strings into usable Arguments
type ArgumentList struct {
	PosArgs   ArgSlice
	NamedArgs ArgMap
}

// GetArgv reads os.Args and turns it into an ArgMap.
func GetArgv(src []string) ArgumentList {
	src = src[1:]
	// Get positional arguments
	posArgs := []Argument{}
	var namedArgsStart int
	for i, arg := range src {
		if strings.HasPrefix(arg, "--") {
			namedArgsStart = i
			break
		}
		posArgs = append(posArgs, Argument{true, arg, true})
	}
	// Get named arguments
	namedArgs := make(ArgMap)
	var (
		last    string
		hasLast bool = false
		arg     string
	)
	for i := namedArgsStart; i < len(src); i++ {
		arg = src[i]
		if arg = strings.TrimSpace(arg); arg == "" {
			continue
		}
		if strings.HasPrefix(arg, "--") {
			arg = arg[1:]
			if hasLast {
				if strings.HasPrefix(last, "!") {
					namedArgs[last[1:]] = Argument{false, "", false}
				} else {
					namedArgs[last] = Argument{false, "", true}
				}
				hasLast = false
			} else {
				last = arg[1:]
				hasLast = true
			}
		} else {
			namedArgs[last] = Argument{true, arg, true}
			hasLast = false
		}
	}
	// Combine and return
	return ArgumentList{posArgs, namedArgs}
}
