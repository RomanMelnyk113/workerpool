package workerpool

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Task func() error

type pool struct {
	wg   sync.WaitGroup
	log  logrus.FieldLogger
	ctx  context.Context
	size int

	// main queue to accepting tasks
	taskQueue chan Task

	// queue to process tasks by workers
	workerQueue chan Task

	Stopped  chan struct{}
	Stopping chan struct{}
}

func NewPool(ctx context.Context, log logrus.FieldLogger, size int) *pool {
	return &pool{
		ctx:         ctx,
		log:         log,
		size:        size,
		taskQueue:   make(chan Task),
		workerQueue: make(chan Task),
		Stopped:     make(chan struct{}),
		Stopping:    make(chan struct{}),
	}
}

// Stop finishes accepting new tasks by closing the tasks channel
func (p *pool) Stop() {
	close(p.Stopping)
	close(p.taskQueue)
}

// Execute accepts task and pass it to the taskQueue
func (p *pool) Execute(task Task) {
	if task != nil {
		// avoid task executing when pool is closing
		select {
		default:
		case <-p.Stopping:
			return
		}
		p.taskQueue <- task
	}
}

// dispatch keep processing the tasks and send to the workers queue
func (p *pool) Run() {
	p.initWorkers()

	go func() {
	Loop:
		for {
			select {
			case task, more := <-p.taskQueue:
				if !more {
					break Loop
				}
				// Got a task to do.
				select {
				case p.workerQueue <- task:
					//p.log.Info("push task to workers")
				}
			case <-p.ctx.Done():
				p.Stop()
				break Loop
			}
		}
		// close workers queue and wait for all tasks
		close(p.workerQueue)
		p.wg.Wait()
		close(p.Stopped)
	}()
}

// worker initializing processing queue with tasks
func (p *pool) worker(id int) {
	p.log.Debugf("spawn worker %d", id)
	defer p.wg.Done()
	for {
		task, more := <-p.workerQueue
		if !more {
			p.log.Debugf("kill worker %d", id)
			break
		}
		if err := task(); err != nil {
			// TODO: add failed tasks to the queue to rerun them???
			p.log.Warnf("task failed: %v", err)
		}
	}
}

// dispatch starts all workers
func (p *pool) initWorkers() {
	if p.size < 1 {
		p.size = 1
	}
	for i := 0; i < p.size; i++ {
		idx := i
		p.wg.Add(1)
		go p.worker(idx)
	}
}
