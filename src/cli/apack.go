package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"ximfect/environ"
	"ximfect/pack"
	"ximfect/tool"
	"ximfect/vm"
)

func packEffect(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	effID := strings.ToLower(ctx.Args.PArgs[0])

	ctx.Log.Debug("Loading effect: " + effID)
	effObj, err := vm.LoadAppdataEffect(effID)
	if err != nil {
		return fmt.Errorf(
			"could not find effect: %s", effID)
	}

	outFileName := effID + "-" + effObj.Metadata.Version + ".fx.xpk"

	ctx.Log.Debug("Packaging...")
	path := environ.DataPath("effects", effID)
	raw, err := pack.GetPackedDirectory(effID, path)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Saving to file: " + outFileName)
	file, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	file.Write(raw)

	return nil
}

func packLib(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: lib-id)")
	}

	libID := strings.ToLower(ctx.Args.PArgs[0])

	ctx.Log.Debug("Loading lib: " + libID)
	libObj, err := vm.LoadAppdataLib(libID)
	if err != nil {
		return fmt.Errorf(
			"could not find lib: %s", libID)
	}

	outFileName := libID + "-" + libObj.Metadata.Version + ".lib.xpk"

	ctx.Log.Debug("Packaging...")
	path := environ.DataPath("libs", libID)
	raw, err := pack.GetPackedDirectory(libID, path)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Saving to file: " + outFileName)
	file, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	file.Write(raw)

	return nil
}

func packGenerator(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: generator-id)")
	}

	genID := strings.ToLower(ctx.Args.PArgs[0])

	ctx.Log.Debug("Loading generator: " + genID)
	gen, err := vm.LoadAppdataGenerator(genID)
	if err != nil {
		return fmt.Errorf(
			"could not find generator: %s", genID)
	}

	outFileName := genID + "-" + gen.Metadata.Version + ".gen.xpk"

	ctx.Log.Debug("Packaging...")
	path := environ.DataPath("generators", genID)
	raw, err := pack.GetPackedDirectory(genID, path)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Saving to file: " + outFileName)
	file, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	file.Write(raw)

	return nil
}

func unpackEffect(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: package)")
	}

	packageFilename := ctx.Args.PArgs[0]

	ctx.Log.Debug("Reading file: " + packageFilename)
	raw, err := ioutil.ReadFile(packageFilename)
	if err != nil {
		return err
	}
	//fmt.Println(raw)

	ctx.Log.Debug("Parsing package...")
	pkg, err := pack.GetPackage(raw)
	//fmt.Println(pkg)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Unpacking...")
	err = pack.UnpackTo(pkg, environ.DataPath("effects", pkg.Name))
	if err != nil {
		return err
	}

	return nil
}

func unpackLib(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: package)")
	}

	packageFilename := ctx.Args.PArgs[0]

	ctx.Log.Debug("Reading file: " + packageFilename)
	raw, err := ioutil.ReadFile(packageFilename)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Parsing package...")
	pkg, err := pack.GetPackage(raw)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Unpacking...")
	err = pack.UnpackTo(pkg, environ.DataPath("libs", pkg.Name))
	if err != nil {
		return err
	}

	return nil
}

func unpackGenerator(ctx *tool.Context) error {
	if len(ctx.Args.PArgs) < 1 {
		return errors.New("not enough arguments (want: package)")
	}

	packageFilename := ctx.Args.PArgs[0]

	ctx.Log.Debug("Reading file: " + packageFilename)
	raw, err := ioutil.ReadFile(packageFilename)
	if err != nil {
		return err
	}
	//fmt.Println(raw)

	ctx.Log.Debug("Parsing package...")
	pkg, err := pack.GetPackage(raw)
	//fmt.Println(pkg)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Unpacking...")
	err = pack.UnpackTo(pkg, environ.DataPath("generators", pkg.Name))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	return

	MasterTool.ToolLog.Debug("Loading actions from apack...")

	peA := tool.NewAction(
		"pack-effect",
		[]string{"pe"},
		"Creates an effect package.",
		tool.QuickPosArgs("effect-id"),
		packEffect)

	plA := tool.NewAction(
		"pack-lib",
		[]string{"pl"},
		"Creates a lib package.",
		tool.QuickPosArgs("lib-id"),
		packLib)

	pgA := tool.NewAction(
		"pack-generator",
		[]string{"pg"},
		"Creates a generator package.",
		tool.QuickPosArgs("generator-id"),
		packGenerator)

	upeA := tool.NewAction(
		"unpack-effect",
		[]string{"upe"},
		"Installs an effect package.",
		tool.QuickPosArgs("effect-id"),
		unpackEffect)

	uplA := tool.NewAction(
		"unpack-lib",
		[]string{"upl"},
		"Installs a lib package.",
		tool.QuickPosArgs("lib-id"),
		unpackLib)

	upgA := tool.NewAction(
		"unpack-generator",
		[]string{"upg"},
		"Installs a generator package.",
		tool.QuickPosArgs("generator-id"),
		unpackGenerator)

	MasterTool.AddManyActions("effects", peA, upeA)
	MasterTool.AddManyActions("libs", plA, uplA)
	MasterTool.AddManyActions("generators", pgA, upgA)
}
