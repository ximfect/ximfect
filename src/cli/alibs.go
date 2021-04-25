package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"ximfect/environ"
	"ximfect/tool"
	"ximfect/vm"
)

func describeLib(ctx *tool.Context, args tool.ArgList) error {
	if len(args.PArgs) < 1 {
		return errors.New("not enough arguments (want: lib-id)")
	}

	// TODO: find a different way to make it case insensitive
	// libName := strings.ToLower(args.PosArgs[0])
	libName := args.PArgs[0]

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

func listLibs(ctx *tool.Context, args tool.ArgList) error {
	var (
		nameFilter   string
		idFilter     string
		authorFilter string
		descFilter   string

		amt = -1
	)

	nameArg, ok := args.NArgs["name"]
	if !ok {
		nameFilter = ""
	} else {
		nameFilter = strings.ToLower(nameArg.Value)
	}

	idArg, ok := args.NArgs["id"]
	if !ok {
		idFilter = ""
	} else {
		idFilter = strings.ToLower(idArg.Value)
	}

	authorArg, ok := args.NArgs["author"]
	if !ok {
		authorFilter = ""
	} else {
		authorFilter = strings.ToLower(authorArg.Value)
	}

	descArg, ok := args.NArgs["desc"]
	if !ok {
		descFilter = ""
	} else {
		descFilter = strings.ToLower(descArg.Value)
	}

	filepath.WalkDir(environ.DataPath("libs"),
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				amt++
				lib, err := vm.LoadAppdataLib(d.Name())
				if err == nil {
					nameRes := strings.Contains(
						strings.ToLower(lib.Metadata.Name), nameFilter)
					idRes := strings.Contains(
						strings.ToLower(lib.Metadata.ID), idFilter)
					authorRes := strings.Contains(
						strings.ToLower(lib.Metadata.Author), authorFilter)
					descRes := strings.Contains(
						strings.ToLower(lib.Metadata.Desc), descFilter)

					if nameRes && idRes && authorRes && descRes {
						fmt.Printf("%s v%s (%s)\n\tby %s\n\t%s\n",
							lib.Metadata.Name,
							lib.Metadata.Version,
							lib.Metadata.ID,
							lib.Metadata.Author,
							lib.Metadata.Desc)
					}
				}
			}
			return nil
		})

	fmt.Println("\nFound", amt, "lib(s)")
	return nil
}

func init() {
	MasterTool.ToolLog.Debug("loading actions from alibs")

	dlA := tool.NewAction(
		"about-lib",
		[]string{"dl"},
		"Shows information about a lib.",
		tool.QuickPosArgs("lib-id"),
		describeLib)

	llA := tool.NewAction(
		"list-libs",
		[]string{"ll"},
		"Shows a list of available libs.",
		tool.QuickNamedArgs(tool.ArgMap{
			"author": tool.QuickArgument(true, "author", false),
			"desc":   tool.QuickArgument(true, "description", false),
			"id":     tool.QuickArgument(true, "effect ID", false),
			"name":   tool.QuickArgument(true, "name", false)}),
		listLibs)

	MasterTool.AddManyActions("libs", dlA, llA)
}
