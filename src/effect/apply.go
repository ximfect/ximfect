/* apllying effects onto images */

package effect

import (
	"fmt"
	"ximfect/libs"

	"github.com/robertkrimen/otto"
	"github.com/ximfect/ximgy"
)

// PrepareVM adds all the API functions to a VM.
func PrepareVM(vm *otto.Otto, img *ximgy.Image) error {
	size := img.Size
	var err error
	err = vm.Set("Require", func(call otto.FunctionCall) otto.Value {
		libName, err := call.Argument(0).ToString()
		if err != nil {
			fmt.Println(err)
			return otto.Value{}
		}
		lib, err := libs.LoadFromAppdata(libName)
		if err != nil {
			fmt.Println(err)
			return otto.Value{}
		}
		libs.ApplyLib(vm, lib)
		return otto.Value{}

	})
	if err != nil {
		return err
	}
	err = vm.Set("ImageSize", func(call otto.FunctionCall) otto.Value {
		sizemap := make(map[string]int)
		sizemap["x"] = size.X
		sizemap["y"] = size.Y
		val, _ := vm.ToValue(sizemap)
		return val
	})
	if err != nil {
		return err
	}
	err = vm.Set("ImageAt", func(call otto.FunctionCall) otto.Value {
		colormap := make(map[string]int)
		colormap["r"] = 0
		colormap["g"] = 0
		colormap["b"] = 0
		colormap["a"] = 255
		obj, err := vm.ToValue(colormap)
		if err != nil {
			fmt.Println(err)
			return otto.Value{}
		}
		x64, err := call.Argument(0).ToInteger()
		if err != nil {
			return obj
		}
		y64, err := call.Argument(1).ToInteger()
		if err != nil {
			return obj
		}
		x := int(x64)
		y := int(y64)
		size := img.Size
		if (x < 0 || y < 0) || (x >= size.X || y >= size.Y) {
			return obj
		}
		r, g, b, a := img.At(x, y).RGBA()
		colormap["r"] = int(r >> 8)
		colormap["g"] = int(g >> 8)
		colormap["b"] = int(b >> 8)
		colormap["a"] = int(a >> 8)
		val, err := vm.ToValue(colormap)
		return val
	})
	if err != nil {
		return err
	}
	return nil
}

// Apply runs the given Effect on the given Image with an empty VM.
func Apply(fx *Effect, img *ximgy.Image) {
	vm := otto.New()
	PrepareVM(vm, img)
	fx.Load(vm)
	img.Iterate(fx.Run)
}

// ApplyOn does what Apply does but on an existing VM.
// This assumes that the VM has been prepared already.
func ApplyOn(vm *otto.Otto, fx *Effect, img *ximgy.Image) {
	fx.Load(vm)
	img.Iterate(fx.Run)
}
