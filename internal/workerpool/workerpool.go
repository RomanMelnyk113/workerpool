package workerpool

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Task func() error

type pool struct {
	wg  sync.WaitGroup
	log logrus.FieldLogger

	size int

	taskQueue   chan Task
	workerQueue chan Task
}

func NewPool(ctx context.Context, log logrus.FieldLogger, size int) *pool {
	p := &pool{
		log:         log,
		size:        size,
		taskQueue:   make(chan Task),
		workerQueue: make(chan Task),
	}

	p.initWorkers()
	go p.dispatch(ctx)

	return p
}

// Stop finishes accepting new tasks by closing the tasks channel
func (p *pool) Stop() {
	close(p.taskQueue)
}

// Execute accepts task and pass it to the taskQueue
func (p *pool) Execute(task Task) {
	if task != nil {
		p.taskQueue <- task
	}
}

// worker initializing processing queue with tasks
func (p *pool) worker(id int) {
	p.log.Println("Spawn worker", id)
	for task := range p.workerQueue {
		p.wg.Add(1)
		//p.log.Println("worker", id, "processing job")
		if err := task(); err != nil {
			// TODO: add failed tasks to the queue to rerun them???

			defer p.wg.Done()
			p.log.Warnf("task error: %w", err)
		}
	}
}

// dispatch starts all workers
func (p *pool) initWorkers() {
	p.log.Info("starting workers")
	for i := 0; i < p.size; i++ {
		go p.worker(i)
	}
}

// dispatch keep processing the tasks and send to the workers queue
func (p *pool) dispatch(ctx context.Context) {
	for {
		select {
		case task, more := <-p.taskQueue:
			if !more {
				p.log.Info("cancel tasks processing")
				p.wg.Wait() // wait for all tasks to be finished
				break
			}
			// Got a task to do.
			select {
			case p.workerQueue <- task:
				//p.log.Info("push task to workers")
			}
		case <-ctx.Done():
			p.wg.Wait() // wait for all tasks to be finished
			return
		}
	}
}
