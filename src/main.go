package main

import (
	"errors"
	"os"
	"strings"
	"ximfect/cli"
	"ximfect/environ"
	"ximfect/tool"
)

func main() {
	t := cli.MasterTool

	environ.EnsureAppdata()

	var (
		err  error
		act  string
		args tool.ArgumentList
	)

	if len(os.Args) > 1 {
		if strings.HasSuffix(os.Args[1], ".xpk") {
			t.ToolLog.Debug("is package")
			args = tool.ArgumentList{
				PosArgs:   tool.ArgSlice{os.Args[1]},
				NamedArgs: tool.ArgMap{}}
			if strings.HasSuffix(os.Args[1], ".fx.xpk") {
				act = "unpack-effect"
			} else if strings.HasSuffix(os.Args[1], ".lib.xpk") {
				act = "unpack-lib"
			} else {
				t.ToolLog.Error("unknown package type: " + os.Args[1])
				err = errors.New("unknown package type: " + os.Args[1])
			}
		} else {
			if len(os.Args) > 2 {
				args = tool.GetArgv(os.Args[2:])
			} else {
				args = tool.ArgumentList{
					PosArgs:   []string{},
					NamedArgs: map[string]tool.Argument{}}
			}
			act = os.Args[1]
		}
	} else {
		act = "help"
	}

	if err != nil {
		os.Exit(1)
	}
	if err = t.RunAction(act, args); err != nil {
		os.Exit(1)
	}
}
