package goqueue

import (
	"context"

	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/job"
	"github.com/MohamedAljoke/goqueue/internal/worker"
)

type (
	Job         = job.Job
	Status      = job.Status
	HandlerFunc = func(ctx context.Context, payload map[string]any) error
)

type Queue struct {
	registry *handler.HandlerRegistry
	pool     *worker.WorkerPool
}

func NewQueue() *Queue {
	registry := handler.NewHandlerRegistry()
	workerCount := 5
	pool := worker.NewWorkerPool(workerCount, registry)

	return &Queue{
		registry: registry,
		pool:     pool,
	}
}

func (q *Queue) SubmitJob(ctx context.Context, jobType string, payload map[string]interface{}, maxRetry int) (*Job, error) {
	_, err := q.registry.Get(jobType)
	if err != nil {
		return nil, err
	}

	job := job.NewJob(maxRetry)
	job.Type = jobType
	job.Payload = payload

	q.pool.Submit(job)

	return job, nil
}

func (q *Queue) RegisterHandler(jobType string, handlerFunc HandlerFunc) {
	q.registry.Register(jobType, handlerFunc)
}

func (q *Queue) Start() {
	q.pool.Start()
}
func (q *Queue) Stop() {
	q.pool.Stop()
}
