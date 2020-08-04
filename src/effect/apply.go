package effect

import (
	"fmt"
	"image"

	"github.com/robertkrimen/otto"
)

// PrepareVM adds all the API functions to a VM.
func PrepareVM(vm *otto.Otto, img *image.RGBA) {
	vm.Set("ImageSize", func(call otto.FunctionCall) otto.Value {
		obj := otto.Object{}
		obj.Set("x", img.Bounds().Size().X)
		obj.Set("y", img.Bounds().Size().Y)
		return obj.Value()
	})
	vm.Set("ImageAt", func(call otto.FunctionCall) otto.Value {
		colormap := make(map[string]int)
		colormap["r"] = 0
		colormap["g"] = 0
		colormap["b"] = 0
		colormap["a"] = 65535
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
		size := img.Bounds().Size()
		if (x < 0 || y < 0) || (x >= size.X || y >= size.Y) {
			return obj
		}
		r, g, b, a := img.At(x, y).RGBA()
		colormap["r"] = int(r)
		colormap["g"] = int(g)
		colormap["b"] = int(b)
		colormap["a"] = int(a)
		val, err := vm.ToValue(colormap)
		return val
	})
}

// Apply runs the given Effect on the given Image with an empty VM.
func Apply(fx *Effect, img *image.RGBA) error {
	vm := otto.New()
	PrepareVM(vm, img)
	return fx.Run(vm, img)
}

// ApplyOn does what Apply does but on an existing VM.
// This assumes that the VM has been prepared already.
func ApplyOn(vm *otto.Otto, fx *Effect, img *image.RGBA) error {
	return fx.Run(vm, img)
}
