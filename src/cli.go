package main

import (
	"os"
	"strings"
	"ximfect/cli"
	"ximfect/environ"
)

func main() {
	tool := cli.GetGtool()

	environ.EnsureAppdata()
	tool.Init()

	var err error

	if len(os.Args) == 1 {
		err = tool.RunAction([]string{"", "help"})
	} else if strings.HasSuffix(os.Args[1], ".fx.xpk") {
		err = tool.RunAction([]string{"", "unpack-effect", "--file", os.Args[1]})
	} else if strings.HasSuffix(os.Args[1], ".lib.xpk") {
		err = tool.RunAction([]string{"", "unpack-lib", "--file", os.Args[1]})
	} else if strings.HasSuffix(os.Args[1], ".xfc") {
		err = tool.RunAction([]string{"", "apply-chain", "--file", os.Args[1]})
	} else {
		err = tool.RunAction(os.Args)
	}

	if err != nil {
		tool.ErrorExit("ERROR:", err)
	}
}
