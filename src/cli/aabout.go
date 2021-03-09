package cli

import (
	"errors"
	"fmt"
	"strings"
	"ximfect/vm"
	"ximfect/tool"
)

func aboutTool(ctx *tool.Context) error {
	fmt.Println(ctx.Tool.Version)
	return nil
}

func describe(ctx *tool.Context) error {
	eff, hasEff := ctx.Args.NamedArgs["effect"]
	lib, hasLib := ctx.Args.NamedArgs["lib"]

	if !(hasEff || hasLib) {
		return errors.New(
			"what should be described? use --effect <id> or --lib <id>")
	}

	effName := strings.ToLower(eff.Value)
	libName := strings.ToLower(lib.Value)

	if hasEff {
		ctx.Log.Debug("Loading effect: " + effName)
		fx, err := vm.LoadAppdataEffect(effName)
		if err != nil {
			return err
		}

		meta := fx.Metadata

		fmt.Printf("======== About %s ========\n", effName)
		fmt.Printf("Name:           %s\n", meta.Name)
		fmt.Printf("Version:        %s\n", meta.Version)
		fmt.Printf("Author:         %s\n", meta.Author)
		fmt.Printf("Description:    %s\n", meta.Desc)
		if len(meta.Preload) > 0 {
			fmt.Printf("Preload:         %v\n", strings.Join(meta.Preload, ", "))
		}
	} else if hasLib {
		ctx.Log.Debug("Loading lib: " + libName)
		library, err := vm.LoadAppdataLib(libName)
		if err != nil {
			return err
		}

		meta := library.Metadata

		fmt.Printf("======== About %s ========\n", libName)
		fmt.Printf("Name:           %s\n", meta.Name)
		fmt.Printf("Version:        %s\n", meta.Version)
		fmt.Printf("Author:         %s\n", meta.Author)
		fmt.Printf("Description:    %s\n", meta.Desc)
		fmt.Printf("Files:\n    [%s]\n", strings.Join(library.Files, "; "))
	}

	return nil
}

func dev(ctx *tool.Context) error {
	//panic("hello")
	return nil
}

func init() {
	MasterTool.ToolLog.Debug("Loading actions from aabout...")

	aboutToolAction := &tool.Action{
		aboutTool,
		"Shows version information.",
		tool.ArgumentList{}}

	describeAction := &tool.Action{
		describe,
		"Shows an effect/lib's information.",
		tool.ArgumentList{
			tool.ArgSlice{},
			tool.ArgMap{
				"effect": tool.Argument{true, "id", false},
				"lib":    tool.Argument{true, "id", false}}}}

	devAction := &tool.Action{
		dev,
		"dev",
		tool.ArgumentList{}}

	MasterTool.AddAction("about-tool", aboutToolAction)
	MasterTool.AddAction("describe", describeAction)
	MasterTool.AddAction("dev", devAction)
}
