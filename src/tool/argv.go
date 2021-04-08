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
type ArgSlice []string

// ArgumentList turns a slice of strings into usable Arguments
type ArgumentList struct {
	PosArgs   ArgSlice
	NamedArgs ArgMap
}

// FormatUsage formats this ArgumentList as if it was usage.
func (l ArgumentList) FormatUsage() string {
	out := ""
	for _, a := range l.PosArgs {
		if strings.HasSuffix(a, "?") {
			out += "[" + a[0:len(a)-1] + "] "
		} else {
			out += "<" + a + "> "
		}
	}

	for n, a := range l.NamedArgs {
		if a.IsValue {
			if a.BoolValue {
				out += "<--" + n + " (" + a.Value + ")] "
			} else {
				out += "[--" + n + " (" + a.Value + ")] "
			}
		} else {
			if a.BoolValue {
				out += "<--" + n + "> "
			} else {
				out += "[--" + n + "] "
			}
		}
	}

	return strings.TrimSpace(out)
}

// GetArgv reads the input string slice and turns it into an ArgumentList.
func GetArgv(src []string) ArgumentList {
	src = append(src, "--")
	// Get positional arguments
	posArgs := []string{}
	var namedArgsStart int
	for i, arg := range src {
		if strings.HasPrefix(arg, "--") {
			namedArgsStart = i
			break
		}
		posArgs = append(posArgs, arg)
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
	)
	for i, arg := range src[namedArgsStart:] {
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

func QuickPosArgs(args ...string) ArgumentList {
	return ArgumentList{ArgSlice(args), ArgMap{}}
}

func QuickNamedArgs(args ArgMap) ArgumentList {
	return ArgumentList{ArgSlice{}, args}
}

func QuickArgument(v bool, d string, r bool) Argument {
	return Argument{v, d, r}
}
