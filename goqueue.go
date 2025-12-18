package goqueue

import (
	"context"
	"fmt"

	"github.com/MohamedAljoke/goqueue/internal/entity"
	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/storage"
	"github.com/MohamedAljoke/goqueue/internal/worker"
)

// we would receive the handlers and the payload here
type Queue struct {
	storage  storage.JobStorage
	registry *handler.Registry
	pool     *worker.WorkerPool
}

func NewQueue() *Queue {
	storage := storage.NewMemoryStorage()
	registry := handler.NewRegistry()
	workerCount := 5
	pool := worker.NewWorkerPool(storage, registry, workerCount)

	return &Queue{
		storage:  storage,
		registry: registry,
		pool:     pool,
	}
}

func (q *Queue) SubmitJob(ctx context.Context, jobType string, payload map[string]interface{}, maxRetry int) (*entity.Job, error) {
	_, err := q.registry.GetHandler(jobType)
	if err != nil {
		return nil, err
	}

	job := entity.NewJob(maxRetry)
	job.Type = jobType
	job.Payload = payload

	q.storage.SaveJob(ctx, job)
	q.pool.Submit(job)

	return job, nil
}

func (q *Queue) GetJob(ctx context.Context, jobID string) (*entity.Job, error) {
	job, exists := q.storage.GetJob(ctx, jobID)
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobID)
	}
	return job, nil
}

func (q *Queue) RegisterHandler(jobType string, handlerFunc entity.HandlerFunc) {
	q.registry.RegisterHandler(jobType, handlerFunc)
}

func (q *Queue) Start(workerCount int) {
	q.pool.Start(workerCount)
}
func (q *Queue) Stop() {
	q.pool.Stop()
}
