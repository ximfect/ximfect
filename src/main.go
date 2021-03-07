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

	if strings.HasSuffix(os.Args[1], ".xpk") {
		t.ToolLog.Debug("is package")
		args = tool.ArgumentList{
			PosArgs:   tool.ArgSlice{os.Args[1]},
			NamedArgs: tool.ArgMap{}}
		if strings.HasSuffix(os.Args[1], ".fx.xpk") {
			err = t.RunAction("unpack-effect", args)
		} else if strings.HasSuffix(os.Args[1], ".lib.xpk") {
			err = t.RunAction("unpack-lib", args)
		} else {
			t.ToolLog.Error("unknown package type: " + os.Args[1])
			err = errors.New("unknown package type: " + os.Args[1])
		}
	} else {
		t.ToolLog.Debug("is NOT package")
		args = tool.GetArgv(os.Args[2:])
		if len(os.Args) == 1 {
			act = "help"
		} else {
			act = os.Args[1]
		}
		err = t.RunAction(act, args)
	}

	if err != nil {
		os.Exit(1)
	}
}
