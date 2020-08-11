// +build ignore

package effect

import (
	"image"
	"runtime"

	"github.com/robertkrimen/otto"
)

// WorkManager manages Workers
type WorkManager struct {
	vm      *otto.Otto
	workers []Worker
}

// SpawnWorkers spawns Workers
func (w WorkManager) SpawnWorkers(size image.Point) {
	workerCount := runtime.NumCPU()
	step := int(size.X / workerCount)
	if size.X%2 == 1 {
		w.spawnWorkersInternal(step, workerCount-1)
		lastWorker := NewWorker(w.vm, step*(workerCount-1), step+1)
		w.workers = append(w.workers, lastWorker)
	} else {
		w.spawnWorkersInternal(step, workerCount)
	}
}

func (w WorkManager) spawnWorkersInternal(step, count int) {
	var worker Worker
	for i := 0; i < count; i++ {
		worker = NewWorker(w.vm, step*i, step)
		w.workers = append(w.workers, worker)
	}
}
