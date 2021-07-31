package cli

import (
	"fmt"
	"sort"
	"ximfect/tool"
	"ximfect/vm"

	"github.com/ximfect/ximgy"
)

func help(ctx *tool.Context, args tool.ArgList) error {
	if len(args.PArgs) == 0 {
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
		arg := args.PArgs[0]
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

func aboutTool(ctx *tool.Context, args tool.ArgList) error {
	fmt.Println(ctx.Tool.Version, "build", tool.Build)
	return nil
}

func license(ctx *tool.Context, args tool.ArgList) error {
	fmt.Println(ctx.Tool.Name, tool.Version, ":", ctx.Tool.Desc)
	fmt.Println(`
    Copyright (C) 2020-2021  qeaml

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
	`)
	return nil
}

func dev(ctx *tool.Context, args tool.ArgList) error {
	g, e := vm.LoadAppdataGenerator("white")
	if e != nil {
		return e
	}
	i, e := g.Apply(1024, 1024, ctx)
	if e != nil {
		return e
	}
	e = ximgy.Save(i, "gentest.png")
	if e != nil {
		return e
	}
	return nil
}

func init() {
	// this function is ran the moment the application runs.
	// add all actions to MasterTool inside of the init() function

	MasterTool.ToolLog.Debug("Loading actions from ainfo...")

	hA := tool.NewAction(
		"help",
		[]string{"h", "?"},
		"Shows help.",
		tool.QuickPosArgs("action or category?"),
		help)

	vA := tool.NewAction(
		"version",
		[]string{"ver", "v"},
		"Shows version information.",
		tool.QuickPosArgs(),
		aboutTool)

	lA := tool.NewAction(
		"license",
		[]string{"licence", "l"},
		"Shows license information.",
		tool.QuickPosArgs(),
		license)

	dA := tool.NewAction(
		"dev",
		[]string{},
		"(for testing)",
		tool.QuickPosArgs(),
		dev)

	MasterTool.AddAction("misc", dA)
	MasterTool.AddManyActions("info", hA, vA, lA)
}
