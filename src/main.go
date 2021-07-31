package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"ximfect/cli"
	"ximfect/environ"
	"ximfect/tool"
)

func main() {
	// The main CLI for the application
	t := cli.MasterTool
	defer t.Close()

	// Make sure that the appdata directories exist
	err := environ.EnsureAppdata()
	if err != nil {
		t.ToolLog.Error("program data structure does not exist: " + err.Error())
		// permission denied
		if os.IsPermission(err) {
			fmt.Println(
				"ximfect could not create it's data directory, which is",
				"required for ximfect to run.", environ.PermWarn)
			// PermWarn suggests a platform-dependent solution to the problem
		}
		os.Exit(2)
	}

	var (
		act  string
		args tool.ArgList
	)

	// If we have at least 1 argument
	if len(os.Args) > 1 {
		// If that argument is a package
		if strings.HasSuffix(os.Args[1], ".xpk") {
			// debug message
			t.ToolLog.Debug("is package")
			// prepare arguments
			args = tool.ArgList{
				PArgs: tool.ArgSlice{os.Args[1]},
				NArgs: tool.ArgMap{}}
			// If the package is an effect package
			if strings.HasSuffix(os.Args[1], ".fx.xpk") {
				// use unpack-effect
				act = "unpack-effect"
				// If the package is a lib package
			} else if strings.HasSuffix(os.Args[1], ".lib.xpk") {
				// use unpack-lib
				act = "unpack-lib"
				// Otherwise we can't unpack an unknown package type
			} else {
				// error out
				t.ToolLog.Error("unknown package type: " + os.Args[1])
				err = errors.New("unknown package type: " + os.Args[1])
			}
			// It's not
		} else {
			// If there are arguments we can parse
			if len(os.Args) > 2 {
				// Parse them
				args = tool.GetArgv(os.Args[2:])
				// If not
			} else {
				// Just use an empty ArgumentList
				args = tool.ArgList{
					PArgs: []string{},
					NArgs: map[string]tool.Arg{}}
			}
			// the action
			act = os.Args[1]
		}
		// We don't, so we run 'help' by default
	} else {
		act = "help"
	}

	// Error out if something went wrong with the args
	if err != nil {
		os.Exit(1)
	}
	// Run the action and error out if necessary
	if err = t.RunAction(act, args); err != nil {
		os.Exit(1)
	}
}
