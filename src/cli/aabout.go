package cli

import (
	"errors"
	"fmt"
	"strings"
	"ximfect/tool"
	"ximfect/vm"
)

func aboutTool(ctx *tool.Context) error {
	fmt.Println(ctx.Tool.Version)
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

func dev(ctx *tool.Context) error {
	//panic("hello")
	return nil
}

func init() {
	// this function is ran the moment the application runs.
	// add all actions to MasterTool inside of the init() function

	MasterTool.ToolLog.Debug("Loading actions from aabout...")

	aboutToolAction := &tool.Action{
		aboutTool,
		"Shows version information.",
		tool.ArgumentList{},
		[]string{"ver", "v"}}

	describeEffectAction := &tool.Action{
		describeEffect,
		"Shows an effect's information.",
		tool.ArgumentList{
			tool.ArgSlice{"effect-id"},
			tool.ArgMap{}},
		[]string{"de"}}

	describeLibAction := &tool.Action{
		describeLib,
		"Shows a lib's information.",
		tool.ArgumentList{
			tool.ArgSlice{"lib-id"},
			tool.ArgMap{}},
		[]string{"dl"}}

	devAction := &tool.Action{
		dev,
		"dev",
		tool.ArgumentList{},
		[]string{}}

	licenseAction := &tool.Action{
		license,
		"Shows license information.",
		tool.ArgumentList{},
		[]string{"l"}}

	MasterTool.AddAction("version", aboutToolAction)
	MasterTool.AddAction("about-effect", describeEffectAction)
	MasterTool.AddAction("about-lib", describeLibAction)
	MasterTool.AddAction("dev", devAction)
	MasterTool.AddAction("license", licenseAction)
}
