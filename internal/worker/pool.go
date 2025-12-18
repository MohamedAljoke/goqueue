package worker

import (
	"context"
	"sync"

	"github.com/MohamedAljoke/goqueue/internal/entity"
	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/storage"
	"github.com/MohamedAljoke/goqueue/internal/usecase"
)

type WorkerPool struct {
	storage   storage.JobStorage
	registry  *handler.Registry
	processor *usecase.JobProcessor

	jobChan chan *entity.Job

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup
}

func NewWorkerPool(storage storage.JobStorage, registry *handler.Registry, workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	poolChan := make(chan *entity.Job, workerCount*2)

	pool := &WorkerPool{
		storage:   storage,
		registry:  registry,
		processor: usecase.NewJobProcessor(storage),
		jobChan:   poolChan,
		ctx:       ctx,
		cancel:    cancel,
	}

	return pool
}

func (wp *WorkerPool) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

}

func (wp *WorkerPool) Submit(job *entity.Job) {
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
		case job := <-wp.jobChan:
			handler, err := wp.registry.GetHandler(job.Type)
			if err != nil {
				job.MarkFailed(err)
				wp.storage.SaveJob(wp.ctx, job)
				continue
			}
			wp.processor.ProcessJob(wp.ctx, job, handler)
		case <-wp.ctx.Done():
			return
		}
	}
}
