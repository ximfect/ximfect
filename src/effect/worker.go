package effect

import (
	"image"

	"github.com/robertkrimen/otto"
)

// Worker applies an Effect to a part of an image
type Worker struct {
	effect      *Effect
	vm          *otto.Otto
	sliceWidth  int
	sliceHeight int
	sliceOffset int
	manager     *WorkerManager
	IsWorking   bool
}

// NewWorker constructs a Worker
func NewWorker(div, off int, fx *Effect, img *image.RGBA, vm *otto.Otto) *Worker {
	tmp := new(Worker)
	tmp.IsWorking = false
	size := img.Bounds().Size()
	step := int(size.X / div)
	tmp.SetEffect(fx)
	tmp.SetVM(vm)
	tmp.SetSlice(step*off, step, size.Y)
	return tmp
}

// SetEffect sets the effect the Worker will use
func (w *Worker) SetEffect(fx *Effect) {
	w.effect = fx
}

// SetVM sets the VM the Worker will use
func (w *Worker) SetVM(vm *otto.Otto) {
	w.vm = vm
}

// SetSlice sets the "slice" of the image the Worker will work on
func (w *Worker) SetSlice(offset, width, height int) {
	w.sliceOffset = offset
	w.sliceWidth = width
	w.sliceHeight = height
}

// Work applies the Worker's effect to the given RGBA image
func (w *Worker) Work(img *image.RGBA, workerID int) error {
	w.IsWorking = true
	var (
		xMin int = w.sliceOffset
		xMax int = xMin + w.sliceWidth
	)
	for x := xMin; x < xMax; x++ {
		for y := 0; y < w.sliceHeight; y++ {
			color, _ := w.effect.Run(x, y, w.vm, img)
			w.manager.AddToQueue(color, workerID)
		}
	}
	w.IsWorking = false
	return nil
}
