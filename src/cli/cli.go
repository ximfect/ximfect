package cli

import (
	"os"

	"ximfect/tool"
)

// MasterTool is the main Tool for this application
var MasterTool = tool.NewTool(
	"ximfect",
	"An effect-based image processing tool.",
	tool.Version)

func init() {
	for _, a := range os.Args {
		if a == "--debug" {
			MasterTool.MasterLog.SetLevel(0)
		}
	}
}
