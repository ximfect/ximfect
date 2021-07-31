/* os.Args parsing */

package tool

import (
	"strings"
)

// Arg represents a CLI argument.
type Arg struct {
	IsValue   bool
	Value     string
	BoolValue bool
}

// ArgMap is a string->Argument map.
type ArgMap map[string]Arg

// ArgSlice is an Argument slice.
type ArgSlice []string

// ArgList turns a slice of strings into usable Arguments
type ArgList struct {
	PArgs ArgSlice
	NArgs ArgMap
}

// FormatUsage formats this ArgumentList as if it was usage.
func (l ArgList) FormatUsage() string {
	out := ""
	for _, a := range l.PArgs {
		if strings.HasSuffix(a, "?") {
			out += "[" + a[0:len(a)-1] + "] "
		} else {
			out += "<" + a + "> "
		}
	}

	for n, a := range l.NArgs {
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
func GetArgv(src []string) ArgList {
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
				namedArgs[name] = Arg{false, "", valueB}
				i--
			} else {
				namedArgs[name] = Arg{true, valueS, true}
			}
		}
	}
	// Combine and return
	return ArgList{posArgs, namedArgs}
}

func QuickPosArgs(args ...string) ArgList {
	return ArgList{ArgSlice(args), ArgMap{}}
}

func QuickNamedArgs(args ArgMap) ArgList {
	return ArgList{ArgSlice{}, args}
}

func QuickArgument(v bool, d string, r bool) Arg {
	return Arg{v, d, r}
}
