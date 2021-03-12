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

func describeEffect(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	effName := strings.ToLower(ctx.Args.PosArgs[0])

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
		fmt.Printf("Preload:        [%v]\n", strings.Join(meta.Preload, ", "))
	}

	return nil
}

func describeLib(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: lib-id)")
	}

	libName := strings.ToLower(ctx.Args.PosArgs[0])

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

	MasterTool.AddAction("version", aboutToolAction)
	MasterTool.AddAction("about-effect", describeEffectAction)
	MasterTool.AddAction("about-lib", describeLibAction)
	MasterTool.AddAction("dev", devAction)
}
