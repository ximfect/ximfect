package cli

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"ximfect/environ"
	"ximfect/fxchain"
	"ximfect/tool"
	"ximfect/vm"

	"github.com/ximfect/ximgy"
)

func applyEffect(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 3 {
		return errors.New("not enough arguments (want: image, effect-id, output)")
	}

	// arguments
	imageFilename := ctx.Args.PosArgs[0]
	effID := ctx.Args.PosArgs[1]
	outputFilename := ctx.Args.PosArgs[2]

	ctx.Log.Debug("Loading effect: " + effID)
	// load effect from appdata
	fx, err := vm.LoadAppdataEffect(effID)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Opening file: " + imageFilename)
	// open the input file with ximgy
	image, err := ximgy.Open(imageFilename)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Applying effect: " + effID)
	// apply the effect
	err = fx.Apply(image, ctx)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Saving output file: " + outputFilename)
	// save the output with ximgy
	err = ximgy.Save(image, outputFilename)
	if err != nil {
		return err
	}

	return nil
}

const (
	// templates for empty effects
	scriptTemplate = "function effect(pixel)\n    -- your code here\n    return {r=pixel[\"r\"], g=pixel[\"g\"], b=pixel[\"b\"], a=pixel[\"a\"]}\nend\n"
	metaTemplate   = "name: Empty Effect\nversion: 1.0.0\nauthor: unknown <>\ndesc: ximfect generated empty effect\n"
)

func initEffect(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	// arguments
	var noTemplate bool
	noTemplateArg, hasNoTemplate := ctx.Args.NamedArgs["no-template"]
	if hasNoTemplate {
		noTemplate = noTemplateArg.BoolValue
	} else {
		noTemplate = false
	}
	effID := strings.ToLower(ctx.Args.PosArgs[0])
	effPath := environ.DataPath("effects", effID)
	scriptPath := environ.Combine(effPath, "effect.lua")
	metaPath := environ.Combine(effPath, "effect.yml")

	ctx.Log.Debug("Creating effect structure")
	// make effect folder
	err := os.Mkdir(effPath, os.ModePerm)
	if err != nil {
		return err
	}
	// make script file
	script, err := os.Create(scriptPath)
	if err != nil {
		return err
	}
	// make meta file
	meta, err := os.Create(metaPath)
	if err != nil {
		return err
	}

	// if the --no-template flag was NOT specified
	if !noTemplate {
		// write the templates to the script and meta files
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

	// tell the user where the effect is
	fmt.Println("View your effect in:", environ.DataPath("effects", effID))
	return nil
}

func applyChain(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 3 {
		return errors.New("not enough arguments (want: image, fx-chain, output)")
	}

	// arguments
	imageFilename := ctx.Args.PosArgs[0]
	chainFilename := ctx.Args.PosArgs[1]
	outputFilename := ctx.Args.PosArgs[2]

	ctx.Log.Debug("Loading FX chain: " + chainFilename)
	// load the fx chain "script"
	raw, err := ioutil.ReadFile(chainFilename)
	if err != nil {
		return err
	}
	src := string(raw)

	ctx.Log.Debug("Loading image: " + imageFilename)
	// load the source image using ximgy
	img, err := ximgy.Open(imageFilename)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Applying FX chain...")
	// apply the fx chain
	res, err := fxchain.Apply(src, img, ctx)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Saving result: " + outputFilename)
	// save the output image using ximgy
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

	// arguments
	var size int
	var err error
	sizeArg, hasSize := ctx.Args.NamedArgs["size"]
	if hasSize {
		if !sizeArg.IsValue {
			return errors.New("argument `size` is not a value")
		}
		size, err = strconv.Atoi(sizeArg.Value)
		if err != nil {
			return errors.New("argument `size` is not a number")
		}
	} else {
		size = 1024
	}
	outputFilename := ctx.Args.PosArgs[0]

	ctx.Log.Debug("Generating test image...")
	// make an empty image
	img := ximgy.MakeEmpty(image.Rect(0, 0, size, size))
	// gradient step
	step := size / 256
	// fill the image with a gradient
	img.Iterate(func(pixel ximgy.Pixel) (color.RGBA, error) {
		return color.RGBA{
			uint8(pixel.X / step),
			uint8(((pixel.X / 2) + (pixel.Y / 2)) / step),
			uint8(pixel.Y / step), 255}, nil
	})

	ctx.Log.Debug("Saving output file: " + outputFilename)
	// save the output using ximgy
	err = ximgy.Save(img, outputFilename)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	MasterTool.ToolLog.Debug("Loading actions from aeffects...")

	aeA := tool.NewAction(
		"apply-effect",
		[]string{"ae"},
		"Applies an effect to an image.",
		tool.QuickPosArgs("image", "effect-id", "output"),
		applyEffect)

	acA := tool.NewAction(
		"apply-chain",
		[]string{"ac"},
		"Applies an FX chain to an image.",
		tool.QuickPosArgs("image", "fx-chain", "output"),
		applyChain)

	ieA := tool.NewAction(
		"init-effect",
		[]string{"ie"},
		"Creates a new effect template.",
		tool.QuickPosArgs("effect-id"),
		initEffect)

	giA := tool.NewAction(
		"gen-image",
		[]string{"gi"},
		"Generates an image.",
		tool.QuickNamedArgs(tool.ArgMap{
			"size": tool.QuickArgument(true, "image size", false)}),
		genImage)

	MasterTool.AddManyActions("effects", aeA, acA, ieA)
	MasterTool.AddAction("images", giA)
}
