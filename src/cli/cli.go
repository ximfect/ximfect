package cli

import (
	"ximfect/tool"
)

var gTool *tool.Tool = tool.NewTool(
	"ximfect",
	tool.Version,
	"Learn more at https://ximfect.github.io")

// GetGtool returns the global Tool
func GetGtool() *tool.Tool {
	return gTool
}
