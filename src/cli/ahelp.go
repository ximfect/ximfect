package cli

import (
	"errors"
	"fmt"
	"sort"
	"ximfect/tool"
)

func help(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) == 0 {
		ctx.Tool.PrintTitle()
		fmt.Println("\n" + ctx.Tool.Desc + "\n")
		fmt.Println("Required arguments are surrounded by <>")
		fmt.Print("Optional arguments are surrounded by []\n\n")
		list := ctx.Tool.GetActionList()
		sort.Strings(list)
		for _, name := range list {
			action, exists := ctx.Tool.GetAction(name)
			if !exists {
				continue
			}

			var desc string
			if len(action.Desc) > 70 {
				desc = action.Desc[0:70]
			} else {
				desc = action.Desc
			}

			var nameFinal string
			nameLong := name + " " + action.Usage.FormatUsage()
			if len(nameLong) > 75 {
				nameFinal = nameLong[0:75]
			} else {
				nameFinal = nameLong
			}

			fmt.Println("     " + nameFinal)
			fmt.Println("          " + desc)
		}
	} else {
		target := ctx.Args.PosArgs[0]
		action, exists := ctx.Tool.GetAction(target)
		if !exists {
			return errors.New("unknown action: " + target)
		}

		var nameFinal string
		nameLong := target + " " + action.Usage.FormatUsage()
		if len(nameLong) > 80 {
			nameFinal = nameLong[0:80]
		} else {
			nameFinal = nameLong
		}

		fmt.Println(nameFinal)
		fmt.Println("    " + action.Desc)
	}
	return nil
}

func init() {
	helpAction := &tool.Action{
		help,
		"Shows a list of actions or an action's description.",
		tool.ArgumentList{
			tool.ArgSlice{"action?"},
			tool.ArgMap{}}}

	MasterTool.AddAction("help", helpAction)
}
