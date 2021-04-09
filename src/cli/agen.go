package cli

import (
	"errors"
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"ximfect/environ"
	"ximfect/tool"
	"ximfect/vm"

	"github.com/ximfect/ximgy"
)

func generate(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 2 {
		return errors.New("not enough arguments (want: generator-id, output)")
	}

	var size image.Point
	var err error
	sizeArg, hasSize := ctx.Args.NArgs["size"]
	if hasSize {
		sizeElems := strings.Split(sizeArg.Value, "x")
		if len(sizeElems) != 2 {
			return fmt.Errorf("invalid resolution format: %s", sizeArg.Value)
		}

		size.X, err = strconv.Atoi(sizeElems[0])
		if err != nil {
			return err
		}

		size.Y, err = strconv.Atoi(sizeElems[1])
		if err != nil {
			return err
		}
	} else {
		size.X = 1024
		size.Y = 1024
	}

	generator, err := vm.LoadAppdataGenerator(ctx.Args.PArgs[0])
	if err != nil {
		return err
	}

	out, err := generator.Apply(size.X, size.Y, ctx)
	if err != nil {
		return err
	}

	err = ximgy.Save(out, ctx.Args.PArgs[1])
	if err != nil {
		return err
	}

	return nil
}

func describeGenerator(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: generator-id)")
	}

	// TODO: find a different way to make it case insensitive
	// effName := strings.ToLower(ctx.Args.PosArgs[0])
	genName := ctx.Args.PArgs[0]

	ctx.Log.Debug("Loading generator: " + genName)
	// load the effect by id from appdata
	fx, err := vm.LoadAppdataGenerator(genName)
	if err != nil {
		return err
	}

	// get the effect's metadata and print it out
	meta := fx.Metadata
	fmt.Printf("======== About %s ========\n", genName)
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

func listGenerators(ctx *tool.Context) error {
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

	filepath.WalkDir(environ.DataPath("generators"),
		func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}
			if d.IsDir() {
				amt++
				gen, err := vm.LoadAppdataGenerator(d.Name())
				if err == nil {
					nameRes := strings.Contains(
						strings.ToLower(gen.Metadata.Name), nameFilter)
					idRes := strings.Contains(
						strings.ToLower(gen.Metadata.ID), idFilter)
					authorRes := strings.Contains(
						strings.ToLower(gen.Metadata.Author), authorFilter)
					descRes := strings.Contains(
						strings.ToLower(gen.Metadata.Desc), descFilter)

					if nameRes && idRes && authorRes && descRes {
						fmt.Printf("%s v%s (%s)\n\tby %s\n\t%s\n",
							gen.Metadata.Name,
							gen.Metadata.Version,
							gen.Metadata.ID,
							gen.Metadata.Author,
							gen.Metadata.Desc)
					}
				}
			}
			return nil
		})

	fmt.Println("\nFound", amt, "generator(s)")
	return nil
}

func init() {
	MasterTool.ToolLog.Debug("loading actions from agen")

	gA := tool.NewAction(
		"generate",
		[]string{"gi"},
		"Uses a generator to create an image.",
		tool.ArgumentList{
			PArgs: tool.ArgSlice{"generator-id", "output"},
			NArgs: tool.ArgMap{
				"size": tool.QuickArgument(true, "WIDTHxHEIGHT", false)}},
		generate)

	dgA := tool.NewAction(
		"about-generator",
		[]string{"dg"},
		"Shows information about a generator.",
		tool.QuickPosArgs("generator-id"),
		describeGenerator)

	lgA := tool.NewAction(
		"list-generators",
		[]string{"lg"},
		"Shows a list of available generators.",
		tool.QuickNamedArgs(tool.ArgMap{
			"author": tool.QuickArgument(true, "author", false),
			"desc":   tool.QuickArgument(true, "description", false),
			"id":     tool.QuickArgument(true, "effect ID", false),
			"name":   tool.QuickArgument(true, "name", false)}),
		listGenerators)

	MasterTool.AddManyActions("generators", gA, dgA, lgA)
}
