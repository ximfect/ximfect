package effect

import (
	"image"
	"image/color"

	"github.com/robertkrimen/otto"
)

// WorkerManager keeps track of Workers and allows to run them in parallel
type WorkerManager struct {
	WorkerAmt int
	effect    *Effect
	vm        *otto.Otto
	queue     [][]QueueEntry
}

// QueueEntry represents an in the WorkerManager queue
type QueueEntry struct {
	X     int
	Y     int
	Color *color.RGBA
}

// NewWorkerManager constructs a new WorkerManager
func NewWorkerManager(workerAmt int, fx *Effect, vm *otto.Otto) *WorkerManager {
	tmp := new(WorkerManager)
	tmp.WorkerAmt = workerAmt
	tmp.SetVM(vm)
	tmp.SetEffect(fx)
	tmp.ClearQueue()
	return tmp
}

// SetVM sets the VM the Manager will use
func (m *WorkerManager) SetVM(vm *otto.Otto) {
	m.vm = vm
}

// SetEffect sets the Effect the Manager will use
func (m *WorkerManager) SetEffect(fx *Effect) {
	m.effect = fx
}

// ClearQueue clears the pixel change queue
func (m *WorkerManager) ClearQueue() {
	m.queue = make([][]QueueEntry, m.WorkerAmt)
	for i := 0; i < m.WorkerAmt; i++ {
		m.queue[i] = []QueueEntry{}
	}
}

// Run creates WorkerAmt Workers and processes the image with in parallel
func (m *WorkerManager) Run(img *image.RGBA) {
	workers := make([]*Worker, m.WorkerAmt)
	for i := 0; i < m.WorkerAmt; i++ {
		workers[i] = NewWorker(m.WorkerAmt, i, m.effect, img, m.vm)
	}
	var isAnyWorking = true
	for i := 0; i < m.WorkerAmt; i++ {
		worker := workers[i]
		go worker.Work(img, i)
	}
	go m.HandleQueue(img, &workers)
	for isAnyWorking {
		isAnyWorking = false
		for _, worker := range workers {
			if worker.IsWorking {
				isAnyWorking = true
			}
		}
	}
	m.ClearQueue()
}

// HandleQueue loops while atleast 1 Worker is working, apply QueueEntry changes one at a time
func (m *WorkerManager) HandleQueue(img *image.RGBA, workers *[]*Worker) {
	isAnyWorking := true
	for isAnyWorking {
		for i := 0; i < len(m.queue); i++ {
			workerQueue := m.queue[i]
			if len(workerQueue) > 0 {
				entry := workerQueue[0]
				img.SetRGBA(entry.X, entry.Y, *entry.Color)
				m.queue[i] = workerQueue[1:]
			}
		}

		isAnyWorking = false
		for _, worker := range *workers {
			if worker.IsWorking {
				isAnyWorking = true
			}
		}
	}
}

// AddToQueue adds a pending color change to a pixel
func (m *WorkerManager) AddToQueue(entry *QueueEntry, id int) {
	m.queue[id] = append(m.queue[id], *entry)
}
