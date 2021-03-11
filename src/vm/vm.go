package vm

import (
	"errors"
	"fmt"
	"math/rand"
	"ximfect/tool"
	"ximfect/environ"

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
		out.RawSetString("r", lua.LNumber(pixel.R))
		out.RawSetString("g", lua.LNumber(pixel.G))
		out.RawSetString("b", lua.LNumber(pixel.B))
		out.RawSetString("a", lua.LNumber(pixel.A))

		L.Push(out)
		return 1
	})
}

func vmSize(img *ximgy.Image) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		out := L.CreateTable(2, 1)
		out.RawSetString("x", lua.LNumber(img.Size.X))
		out.RawSetString("y", lua.LNumber(img.Size.Y))
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

func vmRandom(L *lua.LState) int {
	n := rand.Float32()
	nV := lua.LNumber(n)
	L.Push(nV)
	return 1
}

func vmRandInt(L *lua.LState) int {
	endRaw := L.Get(1)
	if endRaw.Type() != lua.LTNumber {
		return 0
	}
	end := int(endRaw.(lua.LNumber))
	L.Push(lua.LNumber(rand.Intn(end)))
	return 1
}

func vmInspect(L *lua.LState) int {
	val := L.Get(1)
	fmt.Println(val)
	if val.Type() == lua.LTTable {
		t := val.(*lua.LTable)
		fmt.Println(*t)
	}
	return 0
}

func (e *Effect) vm(img *ximgy.Image, ctx *tool.Context) (*lua.LState, error) {
	log := ctx.Log.Sub("VM")

	log.Debug("Creating state...")
	vm := lua.NewState()

	log.Debug("Adding effect...")
	err := vm.DoFile(environ.Combine(e.Dir, "effect.lua"))
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
	vm.SetGlobal("random", vm.NewFunction(vmRandom))
	vm.SetGlobal("randint", vm.NewFunction(vmRandInt))
	vm.SetGlobal("inspect", vm.NewFunction(vmInspect))

	if len(e.Metadata.Preload) > 0 {
		log.Debug("Applying preload...")
		for _, file := range e.Metadata.Preload {
			err = vm.DoFile(environ.Combine(e.Dir, file))
			if err != nil {
				return nil, err
			}
		}
	}
	
	log.Debug("Done!")
	return vm, nil
}