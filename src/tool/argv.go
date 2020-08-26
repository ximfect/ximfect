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
	src = append(src, "--")
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
		name     string
		hasName  bool
		valueS   string
		valueB   bool
		valueIsB bool
		hasValue bool
		arg      string
	)
	for i := namedArgsStart; i < len(src); i++ {
		arg = src[i]
		if !hasValue {
			if strings.HasPrefix(arg, "--") {
				if hasName {
					hasValue = true
					valueIsB = true
					valueB = !strings.HasPrefix(name, "!")
					if !valueB {
						name = name[1:]
					}
				} else if !hasName {
					hasName = true
					name = strings.ToLower(arg[2:])
				}
			} else {
				if hasName {
					hasValue = true
					valueS = arg
				}
			}
		}
		if hasName && hasValue {
			hasName = false
			hasValue = false
			if valueIsB {
				namedArgs[name] = Argument{false, "", valueB}
				i--
			} else {
				namedArgs[name] = Argument{true, valueS, true}
			}
		}
	}
	// Combine and return
	return ArgumentList{posArgs, namedArgs}
}
