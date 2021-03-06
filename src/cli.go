package main

import (
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
		act  string
		args tool.ArgumentList
	)

	if len(os.Args) == 1 {
		act = "help"
	} else {
		act = os.Args[1]
	}

	if strings.HasSuffix(act, ".fx.xpk") {
		args = tool.ArgumentList{
			tool.ArgSlice{},
			tool.ArgMap{
				"file": tool.Argument{true, act, true}}}
		t.RunAction("unpack-effect", args)
	} else if strings.HasSuffix(os.Args[1], ".lib.xpk") {
		args = tool.ArgumentList{
			tool.ArgSlice{},
			tool.ArgMap{
				"file": tool.Argument{true, act, true}}}
		t.RunAction("unpack-lib", args)
	} else {
		args = tool.GetArgv(os.Args[2:])
		t.RunAction(act, args)
	}
}
