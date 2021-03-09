package vm

import (
	"errors"

	"ximfect/tool"

	lua "github.com/yuin/gopher-lua"

	"github.com/ximfect/ximgy"
)

func vmArg(ctx *tool.Context) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		nameLua := L.Get(1)
		if nameLua.Type() != lua.LTString {
			return 0
		}
		name := string(nameLua.(lua.LString))
		
		if val, ok := ctx.Args.NamedArgs[name]; ok {
			if val.IsValue {
				L.Push(lua.LString(val.Value))
			} else {
				L.Push(lua.LBool(val.BoolValue))
			}
			return 1
		}
		return 0
	})
}

func vmAt(img *ximgy.Image) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		xRaw := L.Get(1)
		if xRaw.Type() != lua.LTNumber {
			return 0
		}
		yRaw := L.Get(2)
		if yRaw.Type() != lua.LTNumber {
			return 0
		}

		x := int(xRaw.(lua.LNumber))
		y := int(yRaw.(lua.LNumber))

		pixel := img.At(x, y)
		out := L.CreateTable(4, 1)
		out.RawSet(lua.LString("r"), lua.LNumber(pixel.R))
		out.RawSet(lua.LString("g"), lua.LNumber(pixel.G))
		out.RawSet(lua.LString("b"), lua.LNumber(pixel.B))
		out.RawSet(lua.LString("a"), lua.LNumber(pixel.A))

		L.Push(out)
		return 1
	})
}

func vmSize(img *ximgy.Image) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		out := L.CreateTable(2, 1)
		out.RawSet(lua.LString("x"), lua.LNumber(img.Size.X))
		out.RawSet(lua.LString("y"), lua.LNumber(img.Size.Y))
		L.Push(out)
		return 1
	})
}

func vmImport(L *lua.LState) int {
	var (
		libName lua.LValue
		lib *Lib
		err error
	)
	libName = L.Get(1)
	if libName.Type() != lua.LTString {
		goto end
	}
	lib, err = LoadAppdataLib(string(libName.(lua.LString)))
	if err != nil {
		goto end
	}
	_ = lib.Apply(L)
end:
	return 0
}

func (e *Effect) vm(img *ximgy.Image, ctx *tool.Context) (*lua.LState, error) {
	log := ctx.Log.Sub("VM")

	log.Debug("Creating state...")
	vm := lua.NewState()

	log.Debug("Adding effect...")
	err := vm.DoString(e.source)
	if err != nil {
		return nil, err
	}

	log.Debug("Checking if effect is correct...")
	fxfnVal := vm.GetGlobal("effect")
	if fxfnVal.Type() != lua.LTFunction {
		return nil, errors.New("effect does not define effect() function")
	}

	log.Debug("Adding API functions...")
	vm.SetGlobal("arg", vm.NewFunction(vmArg(ctx)))
	vm.SetGlobal("at", vm.NewFunction(vmAt(img)))
	vm.SetGlobal("size", vm.NewFunction(vmSize(img)))
	vm.SetGlobal("import", vm.NewFunction(vmImport))
	
	log.Debug("Done!")
	return vm, nil
}