package vm

import (
	"errors"
	"fmt"
	"image"
	"math/rand"
	"strconv"
	"ximfect/environ"
	"ximfect/tool"

	lua "github.com/yuin/gopher-lua"

	"github.com/ximfect/ximgy"
)

// gets a named argument from the CLI call
func vmArg(ctx *tool.Context) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		// get first argument from stack
		nameLua := L.Get(1)
		// make sure argument is a string
		if nameLua.Type() != lua.LTString {
			return 0
		}
		// cast to string
		name := string(nameLua.(lua.LString))

		// check if the requested argument is in named arguments
		if val, ok := ctx.Args.NArgs[name]; ok {
			// check if it's a value or a boolean
			if val.IsValue {
				// value; return a string
				L.Push(lua.LString(val.Value))
			} else {
				// boolean; return a boolean
				L.Push(lua.LBool(val.BoolValue))
			}
			// amount of return values (1)
			return 1
		}
		// nothing found, so there are 0 return values
		return 0
	})
}

// returns a pixel from the image at the specified coordinates
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

// returns the size of the image
func vmSize(size image.Point) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		out := L.CreateTable(2, 1)
		out.RawSetString("x", lua.LNumber(size.X))
		out.RawSetString("y", lua.LNumber(size.Y))
		L.Push(out)
		return 1
	})
}

// imports a lib
func vmImport(L *lua.LState) int {
	var (
		libName lua.LValue
		lib     *Lib
		err     error
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

// returns a random float (0.0 to 1.0)
func vmRandom(L *lua.LState) int {
	n := rand.Float32()
	nV := lua.LNumber(n)
	L.Push(nV)
	return 1
}

// returns a random integer from 0 to n
func vmRandInt(L *lua.LState) int {
	endRaw := L.Get(1)
	if endRaw.Type() != lua.LTNumber {
		return 0
	}
	end := int(endRaw.(lua.LNumber))
	L.Push(lua.LNumber(rand.Intn(end)))
	return 1
}

// prints information about the given object
func vmInspect(L *lua.LState) int {
	val := L.Get(1)
	fmt.Println(val)
	if val.Type() == lua.LTTable {
		t := val.(*lua.LTable)
		fmt.Println(*t)
	}
	return 0
}

// turns the given object into an integer (if possible)
func vmInt(L *lua.LState) int {
	val := L.Get(1)
	switch val.Type() {
	case lua.LTNumber:
		L.Push(lua.LNumber(int(val.(lua.LNumber))))
		return 1
	case lua.LTString:
		src := val.(lua.LString).String()
		num, err := strconv.Atoi(src)
		if err != nil {
			return 0
		}
		L.Push(lua.LNumber(num))
		return 1
	default:
		return 0
	}
}

// sends a debug-level message to the logger
func vmDebug(log *tool.Log) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		val := L.Get(1)
		if val.Type() != lua.LTString {
			return 0
		}
		log.Debug(val.String())
		return 0
	})
}

// sends a info-level message to the logger
func vmInfo(log *tool.Log) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		val := L.Get(1)
		if val.Type() != lua.LTString {
			return 0
		}
		log.Info(val.String())
		return 0
	})
}

// sends a warn-level message to the logger
func vmWarn(log *tool.Log) func(L *lua.LState) int {
	return (func(L *lua.LState) int {
		val := L.Get(1)
		if val.Type() != lua.LTString {
			return 0
		}
		log.Warn(val.String())
		return 0
	})
}

// returns a fresh new vm state to be used
func (e *Effect) vm(img *ximgy.Image, ctx *tool.Context) (*lua.LState, error) {
	log := ctx.Log.Sub("VM")

	// create an empty state
	log.Debug("Creating state...")
	vm := lua.NewState()

	// run the effect file
	log.Debug("Adding effect...")
	err := vm.DoFile(environ.Combine(e.Dir, "effect.lua"))
	if err != nil {
		return nil, err
	}

	// check if an effect function is defined
	log.Debug("Checking if effect is correct...")
	fxfnVal := vm.GetGlobal("effect")
	if fxfnVal.Type() != lua.LTFunction {
		return nil, errors.New("effect does not define effect() function")
	}

	// add our public api:
	// to see what each function does, scroll up to it's definition
	log.Debug("Adding API functions...")
	vm.SetGlobal("arg", vm.NewFunction(vmArg(ctx)))
	vm.SetGlobal("at", vm.NewFunction(vmAt(img)))
	vm.SetGlobal("size", vm.NewFunction(vmSize(img.Size)))
	vm.SetGlobal("import", vm.NewFunction(vmImport))
	vm.SetGlobal("random", vm.NewFunction(vmRandom))
	vm.SetGlobal("randint", vm.NewFunction(vmRandInt))
	vm.SetGlobal("inspect", vm.NewFunction(vmInspect))
	vm.SetGlobal("int", vm.NewFunction(vmInt))
	vm.SetGlobal("debug", vm.NewFunction(vmDebug(log)))
	vm.SetGlobal("info", vm.NewFunction(vmInfo(log)))
	vm.SetGlobal("warn", vm.NewFunction(vmWarn(log)))

	// apply preload if necessary
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

// returns a fresh new vm state to be used
func (g *Generator) vm(size image.Point, ctx *tool.Context) (*lua.LState, error) {
	log := ctx.Log.Sub("VM")

	// create an empty state
	log.Debug("Creating state...")
	vm := lua.NewState()

	// run the effect file
	log.Debug("Adding effect...")
	err := vm.DoFile(environ.Combine(g.Dir, "generator.lua"))
	if err != nil {
		return nil, err
	}

	// check if an effect function is defined
	log.Debug("Checking if effect is correct...")
	fxfnVal := vm.GetGlobal("generate")
	if fxfnVal.Type() != lua.LTFunction {
		return nil, errors.New("generator does not define generate() function")
	}

	// add our public api:
	// to see what each function does, scroll up to it's definition
	log.Debug("Adding API functions...")
	vm.SetGlobal("arg", vm.NewFunction(vmArg(ctx)))
	vm.SetGlobal("size", vm.NewFunction(vmSize(size)))
	vm.SetGlobal("import", vm.NewFunction(vmImport))
	vm.SetGlobal("random", vm.NewFunction(vmRandom))
	vm.SetGlobal("randint", vm.NewFunction(vmRandInt))
	vm.SetGlobal("inspect", vm.NewFunction(vmInspect))
	vm.SetGlobal("int", vm.NewFunction(vmInt))
	vm.SetGlobal("debug", vm.NewFunction(vmDebug(log)))
	vm.SetGlobal("info", vm.NewFunction(vmInfo(log)))
	vm.SetGlobal("warn", vm.NewFunction(vmWarn(log)))

	// apply preload if necessary
	if len(g.Metadata.Preload) > 0 {
		log.Debug("Applying preload...")
		for _, file := range g.Metadata.Preload {
			err = vm.DoFile(environ.Combine(g.Dir, file))
			if err != nil {
				return nil, err
			}
		}
	}

	log.Debug("Done!")
	return vm, nil
}
