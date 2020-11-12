package cli

import (
	"errors"
	"fmt"
	"strings"
	"ximfect/effect"
	"ximfect/libs"
	"ximfect/tool"
)

func aboutTool(t *tool.Tool, a tool.ArgumentList) error {
	t.PrintLn(t.GetVersion())
	return nil
}

func describe(t *tool.Tool, a tool.ArgumentList) error {
	eff, hasEff := a.NamedArgs["effect"]
	lib, hasLib := a.NamedArgs["lib"]

	if !(hasEff || hasLib) {
		return errors.New(
			"what should be described? use --effect <id> or --lib <id>")
	}

	effName := strings.ToLower(eff.Value)
	libName := strings.ToLower(lib.Value)

	if hasEff {
		t.VerboseLn("Loading effect:", effName)
		fx, err := effect.LoadFromAppdata(effName)
		if err != nil {
			return err
		}

		meta := fx.Metadata

		t.PrintF("======== About %s ========\n", effName)
		t.PrintF("Name:           %s\n", meta.Name)
		t.PrintF("Version:        %s\n", meta.Version)
		t.PrintF("Author:         %s\n", meta.Author)
		t.PrintF("Description:    %s\n", meta.Desc)
		if len(meta.Preload) > 0 {
			t.PrintF("Preload:         %v\n", strings.Join(meta.Preload, ", "))
		}
	} else if hasLib {
		t.VerboseLn("Loading lib:", libName)
		library, err := libs.LoadFromAppdata(libName)
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

func dev(t *tool.Tool, a tool.ArgumentList) error {
	//panic("hello")
	return nil
}

func init() {
	gTool.VerboseLn("Loading actions from aabout...")
	gTool.AddActionQuick(
		"about-tool",
		"Shows version information",
		"",
		aboutTool)
	gTool.AddActionQuick(
		"describe",
		"Shows an effect/lib's information",
		"[--effect (id)] or [--lib (id)]",
		describe)
	gTool.AddActionQuick(
		"dev",
		"Action for testing purposes",
		"",
		dev)
}
