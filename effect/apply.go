package effect

import (
	"image"

	"github.com/robertkrimen/otto"
)

// Apply runs the given Effect on the given Image with an empty VM.
func Apply(fx *Effect, img *image.RGBA) error {
	size := img.Bounds().Size()
	vm := otto.New()
	vm.Set("IMAGESIZE", [2]int{size.X, size.Y})
	return fx.Run(vm, img)
}

// ApplyOn does what Apply does but on an existing VM.
func ApplyOn(vm *otto.Otto, fx *Effect, img *image.RGBA) error {
	return fx.Run(vm, img)
}
