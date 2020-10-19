/* applying effects onto images */

package effect

import (
	lua "github.com/yuin/gopher-lua"
	"strings"
	"ximfect/libs"
	"ximfect/tool"

	"github.com/ximfect/ximgy"
)

// PrepareVM adds all the API functions to a VM.
func PrepareVM(vm *lua.LState, img *ximgy.Image, args tool.ArgumentList){
	size := img.Size

	// include()
	vm.SetGlobal("include", vm.NewFunction(func(L *lua.LState) int {
		libName := L.ToString(1)
		lib, err := libs.LoadFromAppdata(libName)
		if err != nil {
			return 0
		}
		libs.ApplyLib(vm, lib)
		return 0
	}))
	// size()
	vm.SetGlobal("size", vm.NewFunction(func(L *lua.LState) int {
		out := L.CreateTable(2, 1)
		out.RawSet(lua.LString("x"), lua.LNumber(size.X))
		out.RawSet(lua.LString("y"), lua.LNumber(size.Y))
		L.Push(out)
		return 1
	}))
	// at()
	vm.SetGlobal("at", vm.NewFunction(func(L *lua.LState) int {
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
	}))

	arg := vm.CreateTable(len(args.NamedArgs), 1)
	for name, v := range args.NamedArgs {
		n := strings.ToLower(name)
		if strings.HasPrefix(n, "fx-") {
			if v.IsValue {
				arg.RawSet(lua.LString(n[3:]), lua.LString(v.Value))
			} else {
				arg.RawSet(lua.LString(n[3:]), lua.LNil)
			}
		}
	}
	vm.SetGlobal("args", arg)
}

// Apply runs the given Effect on the given Image with an empty VM.
func Apply(fx *Effect, img *ximgy.Image, tool *tool.Tool, args tool.ArgumentList) error {
	vm := lua.NewState()
	defer vm.Close()
	PrepareVM(vm, img, args)
	fx.Load(vm)
	tool.VerboseLn("- Working...")
	return img.Iterate(fx.Run)
}
