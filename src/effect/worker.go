// +build ignore

package effect

import (
	"image"
	"image/color"

	"github.com/robertkrimen/otto"
)

// Change represents a change in the source image made by a Worker
type Change struct {
	x, y int
	c    color.RGBA
}

// Worker applies an Effect to a slice of an image, storing the changes made
// to be applied later
type Worker struct {
	Changes     []Change
	vm          *otto.Otto
	sliceOffset int
	sliceWidth  int
}

// NewWorker constructs and returns a Worker
func NewWorker(vm *otto.Otto, offset, width int) Worker {
	tmp := new(Worker)
	tmp.SetVM(vm)
	tmp.SetSlice(offset, width)
	return *tmp
}

// SetVM sets the *Worker's VM
func (w *Worker) SetVM(vm *otto.Otto) {
	w.vm = vm
}

// SetSlice sets the *Worker's image slice
func (w *Worker) SetSlice(offset, width int) {
	w.sliceOffset = offset
	w.sliceWidth = width
}

// Work applies the effect to all pixels of the image within the slice
func (w Worker) Work(fx *Effect, img *image.RGBA) {
	height := img.Bounds().Size().Y
	sliceEnd := w.sliceOffset + w.sliceWidth
	var chng Change
	for y := 0; y < height; y++ {
		for x := w.sliceOffset; x < sliceEnd; x++ {
			chng = Change{x, y, img.At(x, y).(color.RGBA)}
			w.Changes = append(w.Changes, chng)
		}
	}
}
