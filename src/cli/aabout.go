package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"ximfect/environ"
	"ximfect/tool"
	"ximfect/vm"
)

func aboutTool(ctx *tool.Context) error {
	fmt.Println(ctx.Tool.Version, "build", tool.Build)
	return nil
}

func license(ctx *tool.Context) error {
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

func describeEffect(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	// TODO: find a different way to make it case insensitive
	// effName := strings.ToLower(ctx.Args.PosArgs[0])
	effName := ctx.Args.PosArgs[0]

	ctx.Log.Debug("Loading effect: " + effName)
	// load the effect by id from appdata
	fx, err := vm.LoadAppdataEffect(effName)
	if err != nil {
		return err
	}

	// get the effect's metadata and print it out
	meta := fx.Metadata
	fmt.Printf("======== About %s ========\n", effName)
	fmt.Printf("Name:           %s\n", meta.Name)
	fmt.Printf("Version:        %s\n", meta.Version)
	fmt.Printf("Author:         %s\n", meta.Author)
	fmt.Printf("Description:    %s\n", meta.Desc)
	// if we have preload, we need to do some extra formatting
	if len(meta.Preload) > 0 {
		fmt.Printf("Preload:        [%v]\n", strings.Join(meta.Preload, ", "))
	}

	return nil
}

func describeLib(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: lib-id)")
	}

	// TODO: find a different way to make it case insensitive
	// libName := strings.ToLower(ctx.Args.PosArgs[0])
	libName := ctx.Args.PosArgs[0]

	ctx.Log.Debug("Loading lib: " + libName)
	// load lib from appdata
	library, err := vm.LoadAppdataLib(libName)
	if err != nil {
		return err
	}

	// get the lib's metadata and print it out
	meta := library.Metadata
	fmt.Printf("======== About %s ========\n", libName)
	fmt.Printf("Name:           %s\n", meta.Name)
	fmt.Printf("Version:        %s\n", meta.Version)
	fmt.Printf("Author:         %s\n", meta.Author)
	fmt.Printf("Description:    %s\n", meta.Desc)
	fmt.Printf("Files:\n    [%s]\n", strings.Join(library.Files, "; "))

	return nil
}

func listEffects(ctx *tool.Context) error {
	var (
		nameFilter   string
		idFilter     string
		authorFilter string
		descFilter   string

		amt = -1
	)

	nameArg, ok := ctx.Args.NamedArgs["name"]
	if !ok {
		nameFilter = ""
	} else {
		nameFilter = strings.ToLower(nameArg.Value)
	}

	idArg, ok := ctx.Args.NamedArgs["id"]
	if !ok {
		idFilter = ""
	} else {
		idFilter = strings.ToLower(idArg.Value)
	}

	authorArg, ok := ctx.Args.NamedArgs["author"]
	if !ok {
		authorFilter = ""
	} else {
		authorFilter = strings.ToLower(authorArg.Value)
	}

	descArg, ok := ctx.Args.NamedArgs["desc"]
	if !ok {
		descFilter = ""
	} else {
		descFilter = strings.ToLower(descArg.Value)
	}

	filepath.WalkDir(environ.DataPath("effects"),
		func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}
			if d.IsDir() {
				amt++
				eff, err := vm.LoadAppdataEffect(d.Name())
				if err == nil {
					nameRes := strings.Contains(
						strings.ToLower(eff.Metadata.Name), nameFilter)
					idRes := strings.Contains(
						strings.ToLower(eff.Metadata.ID), idFilter)
					authorRes := strings.Contains(
						strings.ToLower(eff.Metadata.Author), authorFilter)
					descRes := strings.Contains(
						strings.ToLower(eff.Metadata.Desc), descFilter)

					if nameRes && idRes && authorRes && descRes {
						fmt.Printf("%s v%s (%s)\n\tby %s\n\t%s\n",
							eff.Metadata.Name,
							eff.Metadata.Version,
							eff.Metadata.ID,
							eff.Metadata.Author,
							eff.Metadata.Desc)
					}
				}
			}
			return nil
		})

	fmt.Println("\nFound", amt, "effect(s)")
	return nil
}

func listLibs(ctx *tool.Context) error {
	var (
		nameFilter   string
		idFilter     string
		authorFilter string
		descFilter   string

		amt = -1
	)

	nameArg, ok := ctx.Args.NamedArgs["name"]
	if !ok {
		nameFilter = ""
	} else {
		nameFilter = strings.ToLower(nameArg.Value)
	}

	idArg, ok := ctx.Args.NamedArgs["id"]
	if !ok {
		idFilter = ""
	} else {
		idFilter = strings.ToLower(idArg.Value)
	}

	authorArg, ok := ctx.Args.NamedArgs["author"]
	if !ok {
		authorFilter = ""
	} else {
		authorFilter = strings.ToLower(authorArg.Value)
	}

	descArg, ok := ctx.Args.NamedArgs["desc"]
	if !ok {
		descFilter = ""
	} else {
		descFilter = strings.ToLower(descArg.Value)
	}

	filepath.WalkDir(environ.DataPath("libs"),
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				amt++
				lib, err := vm.LoadAppdataLib(d.Name())
				if err == nil {
					nameRes := strings.Contains(
						strings.ToLower(lib.Metadata.Name), nameFilter)
					idRes := strings.Contains(
						strings.ToLower(lib.Metadata.ID), idFilter)
					authorRes := strings.Contains(
						strings.ToLower(lib.Metadata.Author), authorFilter)
					descRes := strings.Contains(
						strings.ToLower(lib.Metadata.Desc), descFilter)

					if nameRes && idRes && authorRes && descRes {
						fmt.Printf("%s v%s (%s)\n\tby %s\n\t%s\n",
							lib.Metadata.Name,
							lib.Metadata.Version,
							lib.Metadata.ID,
							lib.Metadata.Author,
							lib.Metadata.Desc)
					}
				}
			}
			return nil
		})

	fmt.Println("\nFound", amt, "lib(s)")
	return nil
}

func dev(ctx *tool.Context) error {
	//panic("hello")
	return nil
}

func init() {
	// this function is ran the moment the application runs.
	// add all actions to MasterTool inside of the init() function

	MasterTool.ToolLog.Debug("Loading actions from aabout...")

	deA := tool.NewAction(
		"about-effect",
		[]string{"de"},
		"Shows information about an effect.",
		tool.QuickPosArgs("effect-id"),
		describeEffect)

	dlA := tool.NewAction(
		"about-lib",
		[]string{"dl"},
		"Shows information about a lib.",
		tool.QuickPosArgs("lib-id"),
		describeLib)

	leA := tool.NewAction(
		"list-effects",
		[]string{"le"},
		"Shows a list of available effects.",
		tool.QuickNamedArgs(tool.ArgMap{
			"author": tool.QuickArgument(true, "author", false),
			"desc":   tool.QuickArgument(true, "description", false),
			"id":     tool.QuickArgument(true, "effect ID", false),
			"name":   tool.QuickArgument(true, "name", false)}),
		listEffects)

	llA := tool.NewAction(
		"list-libs",
		[]string{"ll"},
		"Shows a list of available libs.",
		tool.QuickNamedArgs(tool.ArgMap{
			"author": tool.QuickArgument(true, "author", false),
			"desc":   tool.QuickArgument(true, "description", false),
			"id":     tool.QuickArgument(true, "effect ID", false),
			"name":   tool.QuickArgument(true, "name", false)}),
		listLibs)

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
	MasterTool.AddManyActions("info", vA, lA)
	MasterTool.AddManyActions("effects", deA, leA)
	MasterTool.AddManyActions("libs", dlA, llA)
}
