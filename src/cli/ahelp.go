package cli

import (
	"fmt"
	"sort"
	"ximfect/tool"
)

func help(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) == 0 {
		ctx.Tool.PrintTitle()
		fmt.Print("\n" + ctx.Tool.Desc + "\n\n")

		categ := []string{}
		for catName := range ctx.Tool.Categories {
			categ = append(categ, catName)
		}
		sort.Strings(categ)

		for _, catName := range categ {
			cat := ctx.Tool.Categories[catName]
			fmt.Print("    " + catName + "\n")
			fmt.Print("        " + cat.Desc + "\n")
		}
	} else {
		arg := ctx.Args.PosArgs[0]
		cat, isCat := ctx.Tool.Categories[arg]
		act, isAct := ctx.Tool.GetAction(arg)
		if isCat {
			ctx.Tool.PrintTitle()
			fmt.Print("\n" + cat.Desc + "\n\n")
			for _, act := range cat.Actions {
				var desc string
				if len(act.Desc) > 70 {
					desc = act.Desc[0:70]
				} else {
					desc = act.Desc
				}

				var nameFinal string
				nameWithAliases := act.Name
				for _, a := range act.Alias {
					nameWithAliases += "|" + a
				}
				nameLong := nameWithAliases + " " + act.Usage.FormatUsage()
				if len(nameLong) > 75 {
					nameFinal = nameLong[0:75]
				} else {
					nameFinal = nameLong
				}

				fmt.Println("     " + nameFinal)
				fmt.Println("          " + desc)
			}
		} else if isAct {
			fmt.Print("    " + act.Name + "\n")
			fmt.Print("        " + act.Usage.FormatUsage() + "\n\n")
			fmt.Print("        " + act.Desc + "\n")
		} else {
			return fmt.Errorf("i don't know what `%s` is", arg)
		}
	}
	return nil
}

func init() {
	hA := tool.NewAction(
		"help",
		[]string{"h", "?"},
		"Shows help.",
		tool.QuickPosArgs("action or category?"),
		help)

	MasterTool.AddAction("info", hA)
}
