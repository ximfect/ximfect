package cli

import (
	"errors"
	"fmt"
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

func applyEffect(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 3 {
		return errors.New("not enough arguments (want: image, effect-id, output)")
	}

	imageFilename := ctx.Args.PosArgs[0]
	effID := ctx.Args.PosArgs[1]
	outputFilename := ctx.Args.PosArgs[2]

	ctx.Log.Debug("Loading effect: " + effID)
	fx, err := effect.LoadFromAppdata(effID)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Opening file: " + imageFilename)
	image, err := ximgy.Open(imageFilename)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Applying effect: " + effID)
	err = effect.Apply(fx, image, ctx)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Saving output file: " + outputFilename)
	err = ximgy.Save(image, outputFilename)
	if err != nil {
		return err
	}

	return nil
}

const (
	scriptTemplate = "function effect(x, y, pixel) {\n	// write your code here\n	return {r: pixel.r, g: pixel.g, b: pixel.b, a: pixel.a};\n}\n"
	metaTemplate   = "name: Empty Effect\nversion: 1.0.0\nauthor: unknown <>\ndesc: ximfect generated empty effect\n"
)

func initEffect(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	var noTemplate bool
	noTemplateArg, hasNoTemplate := ctx.Args.NamedArgs["no-template"]
	if hasNoTemplate {
		noTemplate = noTemplateArg.BoolValue
	} else {
		noTemplate = false
	}

	effID := strings.ToLower(ctx.Args.PosArgs[0])
	effPath := environ.AppdataPath("effects", effID)
	scriptPath := environ.Combine(effPath, "effect.lua")
	metaPath := environ.Combine(effPath, "effect.yml")

	ctx.Log.Debug("Creating effect structure")
	err := os.Mkdir(effPath, os.ModePerm)
	if err != nil {
		return err
	}
	script, err := os.Create(scriptPath)
	if err != nil {
		return err
	}
	meta, err := os.Create(metaPath)
	if err != nil {
		return err
	}

	if !noTemplate {
		ctx.Log.Debug("Writing file templates...")
		_, err = script.WriteString(scriptTemplate)
		if err != nil {
			return err
		}
		_, err = meta.WriteString(metaTemplate)
		if err != nil {
			return err
		}
	}

	fmt.Println(" -- View your effect in:", environ.AppdataPath("effects", effID))
	return nil
}

func applyChain(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 3 {
		return errors.New("not enough arguments (want: image, fx-chain, output)")
	}

	imageFilename := ctx.Args.PosArgs[0]
	chainFilename := ctx.Args.PosArgs[1]
	outputFilename := ctx.Args.PosArgs[2]

	ctx.Log.Debug("Loading FX chain: " + chainFilename)
	src, err := environ.LoadTextfile(chainFilename)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Loading image: " + imageFilename)
	img, err := ximgy.Open(imageFilename)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Applying FX chain...")
	res, err := fxchain.Apply(src, img, ctx)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Saving result: " + outputFilename)
	err = ximgy.Save(res, outputFilename)
	if err != nil {
		return err
	}

	return nil
}

func genImage(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: output)")
	}

	outputFilename := ctx.Args.PosArgs[0]

	ctx.Log.Debug("Generating test image...")
	amt := 1024
	img := ximgy.MakeEmpty(image.Rect(0, 0, amt, amt))
	step := amt / 256
	img.Iterate(func(pixel ximgy.Pixel) (color.RGBA, error) {
		return color.RGBA{uint8(pixel.X / step), 0, uint8(pixel.Y / step), 255}, nil
	})

	ctx.Log.Debug("Saving output file: " + outputFilename)
	err := ximgy.Save(img, outputFilename)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	MasterTool.ToolLog.Debug("Loading actions from aeffects...")

	applyEffectAction := &tool.Action{
		applyEffect,
		"Applies an effect to an image.",
		tool.ArgumentList{
			tool.ArgSlice{"image", "effect-id", "output"},
			tool.ArgMap{}}}

	applyChainAction := &tool.Action{
		applyChain,
		"Applies an effect chain to an image.",
		tool.ArgumentList{
			tool.ArgSlice{"image", "fx-chain", "output"},
			tool.ArgMap{}}}

	initEffectAction := &tool.Action{
		initEffect,
		"Initializes an empty effect.",
		tool.ArgumentList{
			tool.ArgSlice{"effect-id"},
			tool.ArgMap{
				"no-template": tool.Argument{false, "generate template?", false}}}}

	genImageAction := &tool.Action{
		genImage,
		"Generates an image.",
		tool.ArgumentList{
			tool.ArgSlice{"output"},
			tool.ArgMap{}}}

	MasterTool.AddAction("apply-effect", applyEffectAction)
	MasterTool.AddAction("apply-chain", applyChainAction)
	MasterTool.AddAction("init-effect", initEffectAction)
	MasterTool.AddAction("gen-image", genImageAction)
}
