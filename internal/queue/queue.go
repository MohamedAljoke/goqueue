package queue

import (
	"context"
	"log"

	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/job"
	"github.com/MohamedAljoke/goqueue/internal/storage"
	"github.com/MohamedAljoke/goqueue/internal/worker"
)

// Queue orchestrates job submission and processing
type Queue struct {
	storage  storage.Storage
	handlers *handler.Registry
	pool     *worker.Pool
}

// Config holds queue configuration
type Config struct {
	WorkerCount int
	BufferSize  int
}

// New creates a new queue
func New(cfg Config, store storage.Storage) *Queue {
	handlers := handler.NewRegistry()
	pool := worker.NewPool(cfg.WorkerCount, cfg.BufferSize, store, handlers)

	return &Queue{
		storage:  store,
		handlers: handlers,
		pool:     pool,
	}
}

// RegisterHandler registers a handler for a job type
func (q *Queue) RegisterHandler(jobType string, fn handler.Func) {
	q.handlers.Register(jobType, fn)
	log.Printf("[QUEUE] Registered handler for job type: %s", jobType)
}

// Submit adds a new job to the queue
func (q *Queue) Submit(ctx context.Context, jobType string, payload map[string]interface{}, maxRetry int) (string, error) {
	j := job.New(jobType, payload, maxRetry)

	// Save to storage
	if err := q.storage.Save(ctx, j); err != nil {
		return "", err
	}

	// Enqueue for processing
	q.pool.Enqueue(j)

	log.Printf("[QUEUE] Job %s (%s) submitted", j.ID, jobType)
	return j.ID, nil
}

// GetJob retrieves a job by ID
func (q *Queue) GetJob(ctx context.Context, id string) (*job.Job, error) {
	return q.storage.Get(ctx, id)
}

// Start begins processing jobs
func (q *Queue) Start(ctx context.Context) {
	q.pool.Start(ctx)
}
