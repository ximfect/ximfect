package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"ximfect/effect"
	"ximfect/environ"
	"ximfect/tool"

	"github.com/ximfect/ximgy"
)

var gTool *tool.Tool = tool.NewTool(
	"ximfect",
	"0.2.0",
	"Learn more at https://github.com/qeamlgit/ximfect")

func _version(t *tool.Tool, a tool.ArgumentList) error {
	fmt.Println(t.GetVersion())
	return nil
}

func _apply(t *tool.Tool, a tool.ArgumentList) error {
	eff, hasEff := a.NamedArgs["effect"]
	file, hasFile := a.NamedArgs["file"]
	out, hasOut := a.NamedArgs["out"]

	if !hasEff {
		return errors.New(
			"missing effect argument, specify with --effect <id>")
	}
	if !hasFile {
		return errors.New(
			"missing input file, specify with --file <filename>")
	}
	if !hasOut {
		return errors.New(
			"missing output file, specify with --out <filename>")
	}

	effName := eff.Value
	inFileName := file.Value
	outFileName := out.Value

	t.VerboseLn("[1/4] Loading effect:", effName)
	fx, err := effect.LoadFromAppdata(effName)
	if err != nil {
		return err
	}

	t.VerboseLn("[2/4] Opening file:", inFileName)
	inFile, err := ximgy.Open(inFileName)
	if err != nil {
		return err
	}

	t.VerboseLn("[3/4] Applying effect...")
	effect.Apply(fx, inFile)

	t.VerboseLn("[4/4] Saving output file:", outFileName)
	err = ximgy.Save(inFile, outFileName)
	if err != nil {
		return err
	}

	t.VerboseLn("Finished!")
	return nil
}

func _about(t *tool.Tool, a tool.ArgumentList) error {
	eff, hasEff := a.NamedArgs["effect"]

	if !hasEff {
		return errors.New(
			"missing effect argument, specify with --effect <id>")
	}

	effName := strings.ToLower(eff.Value)

	t.VerboseLn("Loading effect:", effName)
	fx, err := effect.LoadFromAppdata(effName)
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

	return nil
}

func _pack(t *tool.Tool, a tool.ArgumentList) error {
	eff, hasEff := a.NamedArgs["effect"]
	out, hasOut := a.NamedArgs["out"]

	if !hasEff {
		return errors.New(
			"missing effect argument, specify with --effect <id>")
	}
	if !hasOut {
		return errors.New(
			"missing output file, specify with --out <filename>")
	}

	effName := strings.ToLower(eff.Value)
	outFileName := out.Value

	t.VerboseLn("Loading effect:", effName)
	_, err := effect.LoadFromAppdata(effName)
	if err != nil {
		return fmt.Errorf(
			"could not find effect: %s", effName)
	}

	t.VerboseLn("Creating zip archive:", outFileName)
	err = environ.ZipIt(environ.AppdataPath("effects", effName), outFileName)
	if err != nil {
		return err
	}

	t.VerboseLn("Finished!")
	return nil
}

func _unpack(t *tool.Tool, a tool.ArgumentList) error {
	file, hasFile := a.NamedArgs["file"]

	if !hasFile {
		return errors.New(
			"missing input file, specify with --file <filename>")
	}

	inFileName := file.Value

	t.VerboseLn("Unpacking file:", inFileName)
	err := environ.Unzip(inFileName, environ.AppdataPath("effects"))
	if err != nil {
		return err
	}

	t.VerboseLn("Finished!")
	return nil
}

func _test(t *tool.Tool, a tool.ArgumentList) error {
	out, hasOut := a.NamedArgs["out"]

	if !hasOut {
		return errors.New(
			"missing output file, specify with --out <filename>")
	}

	outFileName := out.Value

	t.VerboseLn("Generating test image...")
	amt := 1024
	img := ximgy.MakeEmpty(image.Rect(0, 0, amt, amt))
	step := amt / 256
	img.Iterate(func(pixel ximgy.Pixel) color.RGBA {
		return color.RGBA{uint8(pixel.X / step), uint8(pixel.Y / step), 0, 255}
	})

	t.VerboseLn("Saving output file:", outFileName)
	err := ximgy.Save(img, outFileName)
	if err != nil {
		return err
	}

	t.VerboseLn("Finished!")
	return nil
}

func main() {
	gTool.Init()
	gTool.AddAction("version", _version, "Shows the version")
	gTool.AddAction("apply", _apply, "Applies an effect")
	gTool.AddAction("about", _about, "Shows information about and effect")
	gTool.AddAction("pack", _pack, "Packs an effect into a zip archive")
	gTool.AddAction("unpack", _unpack, "Unpacks and installs an effect")
	gTool.AddAction("test", _test, "Generates a test image")

	var err error

	if len(os.Args) == 1 {
		err = gTool.RunAction([]string{"", "help"})
	} else if strings.HasSuffix(os.Args[1], ".zip") {
		err = gTool.RunAction([]string{"", "unpack", "--file", os.Args[1]})
	} else {
		err = gTool.RunAction(os.Args)
	}

	if err != nil {
		gTool.ErrorExit("ERROR:", err)
	}
}
