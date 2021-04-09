package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"ximfect/environ"
	"ximfect/fxchain"
	"ximfect/tool"
	"ximfect/vm"

	"github.com/ximfect/ximgy"
)

func applyEffect(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 3 {
		return errors.New("not enough arguments (want: image, effect-id, output)")
	}

	// arguments
	imageFilename := ctx.Args.PArgs[0]
	effID := ctx.Args.PArgs[1]
	outputFilename := ctx.Args.PArgs[2]

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
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	// arguments
	var noTemplate bool
	noTemplateArg, hasNoTemplate := ctx.Args.NArgs["no-template"]
	if hasNoTemplate {
		noTemplate = noTemplateArg.BoolValue
	} else {
		noTemplate = false
	}
	effID := strings.ToLower(ctx.Args.PArgs[0])
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
	if len(ctx.Args.PArgs) < 3 {
		return errors.New("not enough arguments (want: image, fx-chain, output)")
	}

	// arguments
	imageFilename := ctx.Args.PArgs[0]
	chainFilename := ctx.Args.PArgs[1]
	outputFilename := ctx.Args.PArgs[2]

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

func describeEffect(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	// TODO: find a different way to make it case insensitive
	// effName := strings.ToLower(ctx.Args.PosArgs[0])
	effName := ctx.Args.PArgs[0]

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

func listEffects(ctx *tool.Context) error {
	var (
		nameFilter   string
		idFilter     string
		authorFilter string
		descFilter   string

		amt = -1
	)

	nameArg, ok := ctx.Args.NArgs["name"]
	if !ok {
		nameFilter = ""
	} else {
		nameFilter = strings.ToLower(nameArg.Value)
	}

	idArg, ok := ctx.Args.NArgs["id"]
	if !ok {
		idFilter = ""
	} else {
		idFilter = strings.ToLower(idArg.Value)
	}

	authorArg, ok := ctx.Args.NArgs["author"]
	if !ok {
		authorFilter = ""
	} else {
		authorFilter = strings.ToLower(authorArg.Value)
	}

	descArg, ok := ctx.Args.NArgs["desc"]
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

func init() {
	MasterTool.ToolLog.Debug("Loading actions from afx...")

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

	deA := tool.NewAction(
		"about-effect",
		[]string{"de"},
		"Shows information about an effect.",
		tool.QuickPosArgs("effect-id"),
		describeEffect)

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

	MasterTool.AddManyActions("effects", aeA, acA, ieA, deA, leA)
}
