package cli

import (
	"errors"
	"image"
	"image/color"
	"os"
	"strings"

	"ximfect/effect"
	"ximfect/environ"
	"ximfect/fxchain"
	"ximfect/tool"

	"github.com/ximfect/ximgy"
)

func applyEffect(t *tool.Tool, a tool.ArgumentList) error {
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

	t.VerboseLn("Loading effect:", effName)
	fx, err := effect.LoadFromAppdata(effName)
	if err != nil {
		return err
	}

	t.VerboseLn("Opening file:", inFileName)
	inFile, err := ximgy.Open(inFileName)
	if err != nil {
		return err
	}

	t.PrintLn("Applying effect:", effName)
	err = effect.Apply(fx, inFile, t, a)
	if err != nil {
		return err
	}

	t.VerboseLn("Saving output file:", outFileName)
	err = ximgy.Save(inFile, outFileName)
	if err != nil {
		return err
	}

	return nil
}

const (
	scriptTemplate string = "function effect(x, y, pixel) {\n	// write your code here\n	return {r: pixel.r, g: pixel.g, b: pixel.b, a: pixel.a};\n}\n"
	metaTemplate   string = "name: Empty Effect\nversion: 1.0.0\nauthor: unknown <>\ndesc: ximfect generated empty effect\n"
)

func initEffect(t *tool.Tool, a tool.ArgumentList) error {
	eff, hasEff := a.NamedArgs["effect"]

	if !hasEff {
		return errors.New(
			"missing effect argument, specify with --effect <id>")
	}

	effName := strings.ToLower(eff.Value)

	t.PrintLn("Creating effect structure")
	err := os.Mkdir(environ.AppdataPath("effects", effName), os.ModePerm)
	if err != nil {
		return err
	}
	script, err := os.Create(environ.AppdataPath("effects", effName, "effect.js"))
	if err != nil {
		return err
	}
	meta, err := os.Create(environ.AppdataPath("effects", effName, "effect.yml"))
	if err != nil {
		return err
	}

	t.VerboseLn("Writing file templates...")
	_, err = script.WriteString(scriptTemplate)
	if err != nil {
		return err
	}
	_, err = meta.WriteString(metaTemplate)
	if err != nil {
		return err
	}

	t.PrintLn(" -- View your effect in:", environ.AppdataPath("effects", effName))
	return nil
}

func applyChain(t *tool.Tool, a tool.ArgumentList) error {
	file, hasFile := a.NamedArgs["file"]
	out, hasOut := a.NamedArgs["out"]
	inp, hasInp := a.NamedArgs["img"]

	if !hasFile {
		return errors.New(
			"missing input file, specify with --file <filename>")
	}
	if !hasOut {
		return errors.New(
			"missing output file, specify with --out <filename>")
	}
	if !hasInp {
		return errors.New(
			"missing input image, specify with --img <filename>")
	}

	inFileName := file.Value
	outFileName := out.Value
	inpFileName := inp.Value

	t.VerboseLn("Loading FX chain: ", inFileName)
	src, err := environ.LoadTextfile(inFileName)
	if err != nil {
		return err
	}

	t.VerboseLn("Loading image:", inpFileName)
	img, err := ximgy.Open(inpFileName)
	if err != nil {
		return err
	}

	t.PrintLn("Applying FX chain...")
	res, err := fxchain.Apply(src, img, t)
	if err != nil {
		return err
	}

	t.VerboseLn("Saving result:", outFileName)
	err = ximgy.Save(res, outFileName)
	if err != nil {
		return err
	}

	return nil
}

func genImage(t *tool.Tool, a tool.ArgumentList) error {
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
	img.Iterate(func(pixel ximgy.Pixel) (color.RGBA, error) {
		return color.RGBA{uint8(pixel.X / step), 0, uint8(pixel.Y / step), 255}, nil
	})

	t.VerboseLn("Saving output file:", outFileName)
	err := ximgy.Save(img, outFileName)
	if err != nil {
		return err
	}

	t.VerboseLn("Finished!")
	return nil
}

func init() {
	gTool.VerboseLn("Loading actions from aeffects...")
	gTool.AddActionQuick(
		"apply-effect",
		"Applies an effect to an image",
		"--effect (id) --file (name) --out (name)",
		applyEffect)
	gTool.AddActionQuick(
		"init-effect",
		"Initializes an empty effect template",
		"--effect (id)",
		initEffect)
	gTool.AddActionQuick(
		"apply-chain",
		"Applies an effect chain to an image",
		"--file (name) --out (name) --img (name)",
		applyChain)
	gTool.AddActionQuick(
		"gen-image",
		"Generates an image",
		"--out (name)",
		genImage)
}
