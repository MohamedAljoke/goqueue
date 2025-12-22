package worker

import (
	"context"
	"sync"
	"time"

	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/job"
)

type WorkerPool struct {
	jobChan  chan *job.Job
	handlers *handler.HandlerRegistry

	workerCount int
	ctx         context.Context
	cancel      context.CancelFunc

	wg sync.WaitGroup
}

func NewWorkerPool(workerCount int, handlers *handler.HandlerRegistry) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	poolChan := make(chan *job.Job, workerCount*2)

	pool := &WorkerPool{
		jobChan:     poolChan,
		ctx:         ctx,
		cancel:      cancel,
		handlers:    handlers,
		workerCount: workerCount,
	}

	return pool
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

}

func (wp *WorkerPool) Submit(job *job.Job) {
	wp.jobChan <- job
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobChan)
	wp.wg.Wait()
}
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case j, ok := <-wp.jobChan:
			if !ok {
				return
			}
			wp.processJob(j)

		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) processJob(j *job.Job) {
	handler, err := wp.handlers.Get(j.Type)
	if err != nil {
		_ = j.MarkFailed(err)

		return
	}

	if err := j.MarkRunning(); err != nil {
		_ = j.MarkFailed(err)

		return
	}

	err = handler(wp.ctx, j.Payload)
	if err == nil {
		_ = j.MarkCompleted()
		return
	}

	_ = j.MarkFailed(err)

	if j.Status == job.StatusPending {
		time.AfterFunc(j.BackoffDuration(), func() {
			wp.Submit(j)
		})
	}
}
