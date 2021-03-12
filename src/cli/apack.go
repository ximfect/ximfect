package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"ximfect/environ"
	"ximfect/vm"
	"ximfect/pack"
	"ximfect/tool"
)

func packEffect(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: effect-id)")
	}

	effID := strings.ToLower(ctx.Args.PosArgs[0])

	ctx.Log.Debug("Loading effect: " + effID)
	effObj, err := vm.LoadAppdataEffect(effID)
	if err != nil {
		return fmt.Errorf(
			"could not find effect: %s", effID)
	}

	outFileName := effID + "-" + effObj.Metadata.Version + ".fx.xpk"

	ctx.Log.Debug("Packaging...")
	path := environ.AppdataPath("effects", effID)
	raw, err := pack.GetPackedDirectory(path)
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
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: lib-id)")
	}

	libID := strings.ToLower(ctx.Args.PosArgs[0])

	ctx.Log.Debug("Loading lib: " + libID)
	libObj, err := vm.LoadAppdataLib(libID)
	if err != nil {
		return fmt.Errorf(
			"could not find lib: %s", libID)
	}

	outFileName := libID + "-" + libObj.Metadata.Version + ".lib.xpk"

	ctx.Log.Debug("Packaging...")
	path := environ.AppdataPath("libs", libID)
	raw, err := pack.GetPackedDirectory(path)
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
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: package)")
	}

	packageFilename := ctx.Args.PosArgs[0]

	ctx.Log.Debug("Reading file: " + packageFilename)
	raw, err := environ.LoadRawfile(packageFilename)
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
	err = pack.UnpackTo(pkg, environ.AppdataPath("effects", pkg.Name))
	if err != nil {
		return err
	}

	return nil
}

func unpackLib(ctx *tool.Context) error {
	if len(ctx.Args.PosArgs) < 1 {
		return errors.New("not enough arguments (want: package)")
	}

	packageFilename := ctx.Args.PosArgs[0]

	ctx.Log.Debug("Reading file: " + packageFilename)
	raw, err := environ.LoadRawfile(packageFilename)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Parsing package...")
	pkg, err := pack.GetPackage(raw)
	if err != nil {
		return err
	}

	ctx.Log.Debug("Unpacking...")
	err = pack.UnpackTo(pkg, environ.AppdataPath("libs", pkg.Name))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	MasterTool.ToolLog.Debug("Loading actions from apack...")

	packEffectAction := &tool.Action{
		packEffect,
		"Packs an effect.",
		tool.ArgumentList{
			tool.ArgSlice{"effect-id"},
			tool.ArgMap{}},
		[]string{"pe"}}

	packLibAction := &tool.Action{
		packLib,
		"Packs a lib.",
		tool.ArgumentList{
			tool.ArgSlice{"lib-id"},
			tool.ArgMap{}},
		[]string{"pl"}}

	unpackEffectAction := &tool.Action{
		unpackEffect,
		"Unpacks an effect.",
		tool.ArgumentList{
			tool.ArgSlice{"package"},
			tool.ArgMap{}},
		[]string{"upe"}}

	unpackLibAction := &tool.Action{
		unpackLib,
		"Unpacks a lib.",
		tool.ArgumentList{
			tool.ArgSlice{"package"},
			tool.ArgMap{}},
		[]string{"upl"}}

	MasterTool.AddAction("pack-effect", packEffectAction)
	MasterTool.AddAction("pack-lib", packLibAction)
	MasterTool.AddAction("unpack-effect", unpackEffectAction)
	MasterTool.AddAction("unpack-lib", unpackLibAction)
}
