package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"ximfect/effect"
	"ximfect/environ"
	"ximfect/libs"
	"ximfect/pack"
	"ximfect/tool"
)

func packEffect(t *tool.Tool, a tool.ArgumentList) error {
	eff, hasEff := a.NamedArgs["effect"]

	if !hasEff {
		return errors.New(
			"missing input argument, specify with --effect <id>")
	}

	effName := strings.ToLower(eff.Value)

	t.VerboseLn("Loading effect:", effName)
	effObj, err := effect.LoadFromAppdata(effName)
	if err != nil {
		return fmt.Errorf(
			"could not find effect: %s", effName)
	}

	outFileName := effName + "-" + effObj.Metadata.Version + ".fx.xpk"

	t.PrintLn("Packaging...")
	path := environ.AppdataPath("effects", effName)
	raw, err := pack.GetPackedDirectory(path)
	if err != nil {
		return err
	}

	t.VerboseLn("Saving to file:", outFileName)
	file, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	file.Write(raw)

	return nil
}

func packLib(t *tool.Tool, a tool.ArgumentList) error {
	lib, hasLib := a.NamedArgs["lib"]

	if !hasLib {
		return errors.New(
			"missing input argument, specify with --lib <id>")
	}

	libName := strings.ToLower(lib.Value)

	t.VerboseLn("Loading lib:", libName)
	libObj, err := libs.LoadFromAppdata(libName)
	if err != nil {
		return fmt.Errorf(
			"could not find lib: %s", libName)
	}

	outFileName := libName + "-" + libObj.Metadata.Version + ".lib.xpk"

	t.PrintLn("Packaging...")
	path := environ.AppdataPath("libs", libName)
	raw, err := pack.GetPackedDirectory(path)
	if err != nil {
		return err
	}

	t.VerboseLn("Saving to file:", outFileName)
	file, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	file.Write(raw)

	return nil
}

func unpackEffect(t *tool.Tool, a tool.ArgumentList) error {
	file, hasFile := a.NamedArgs["file"]

	if !hasFile {
		return errors.New(
			"missing input file, specify with --file <filename>")
	}

	inFileName := file.Value

	t.VerboseLn("Reading file:", inFileName)
	raw, err := environ.LoadRawfile(inFileName)
	if err != nil {
		return err
	}
	//fmt.Println(raw)

	t.VerboseLn("Parsing package...")
	pkg, err := pack.GetPackage(raw)
	//fmt.Println(pkg)
	if err != nil {
		return err
	}

	t.PrintLn("Unpacking...")
	err = pack.UnpackTo(pkg, environ.AppdataPath("effects", pkg.Name))
	if err != nil {
		return err
	}

	return nil
}

func unpackLib(t *tool.Tool, a tool.ArgumentList) error {
	file, hasFile := a.NamedArgs["file"]

	if !hasFile {
		return errors.New(
			"missing input file, specify with --file <filename>")
	}

	inFileName := file.Value

	t.VerboseLn("Reading file:", inFileName)
	raw, err := environ.LoadRawfile(inFileName)
	if err != nil {
		return err
	}

	t.VerboseLn("Parsing package...")
	pkg, err := pack.GetPackage(raw)
	if err != nil {
		return err
	}

	t.PrintLn("Unpacking...")
	err = pack.UnpackTo(pkg, environ.AppdataPath("libs", pkg.Name))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	gTool.VerboseLn("Loading actions from apack...")
	gTool.AddActionQuick(
		"pack-effect",
		"Packs an effect",
		"--effect (id)",
		packEffect)
	gTool.AddActionQuick(
		"pack-lib",
		"Packs a lib",
		"--lib (id)",
		packLib)
	gTool.AddActionQuick(
		"unpack-effect",
		"Unpacks an effect",
		"--effect (id)",
		unpackEffect)
	gTool.AddActionQuick(
		"unpack-lib",
		"Unpacks a lib",
		"--lib (id)",
		unpackLib)
}
