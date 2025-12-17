package goqueue

import (
	"context"

	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/queue"
	"github.com/MohamedAljoke/goqueue/internal/storage"
)

// HandlerFunc is the signature for job handlers
type HandlerFunc func(ctx context.Context, payload map[string]interface{}) error

// Job represents a job (exposed for status checking)
type Job struct {
	ID       string
	Type     string
	Status   string
	Attempts int
	Error    string
}

// Queue is the main entry point for the job processing system
type Queue struct {
	q *queue.Queue
}

// Option is a functional option for configuring the queue
type Option func(*config)

type config struct {
	workerCount int
	bufferSize  int
	storage     storage.Storage
}

// WithWorkers sets the number of concurrent workers
func WithWorkers(count int) Option {
	return func(c *config) {
		c.workerCount = count
	}
}

// WithBufferSize sets the job buffer size
func WithBufferSize(size int) Option {
	return func(c *config) {
		c.bufferSize = size
	}
}

// WithStorage sets a custom storage backend
func WithStorage(store storage.Storage) Option {
	return func(c *config) {
		c.storage = store
	}
}

// New creates a new job queue with the given options
func New(opts ...Option) *Queue {
	// Default configuration
	cfg := &config{
		workerCount: 3,
		bufferSize:  10,
		storage:     storage.NewMemory(),
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Create internal queue
	q := queue.New(queue.Config{
		WorkerCount: cfg.workerCount,
		BufferSize:  cfg.bufferSize,
	}, cfg.storage)

	return &Queue{q: q}
}

// RegisterHandler registers a handler function for a job type
func (q *Queue) RegisterHandler(jobType string, fn HandlerFunc) {
	q.q.RegisterHandler(jobType, handler.Func(fn))
}

// Submit submits a new job to the queue
func (q *Queue) Submit(ctx context.Context, jobType string, payload map[string]interface{}, maxRetry int) (string, error) {
	return q.q.Submit(ctx, jobType, payload, maxRetry)
}

// GetJob retrieves job status by ID
func (q *Queue) GetJob(ctx context.Context, jobID string) (*Job, error) {
	j, err := q.q.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return &Job{
		ID:       j.ID,
		Type:     j.Type,
		Status:   string(j.Status),
		Attempts: j.Attempts,
		Error:    j.Error,
	}, nil
}

// Start begins processing jobs (blocking call)
func (q *Queue) Start(ctx context.Context) {
	q.q.Start(ctx)
}
