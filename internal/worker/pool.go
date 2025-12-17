package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/job"
	"github.com/MohamedAljoke/goqueue/internal/storage"
)

// Pool manages a pool of workers
type Pool struct {
	workerCount int
	jobs        chan *job.Job
	storage     storage.Storage
	handlers    *handler.Registry
}

// NewPool creates a new worker pool
func NewPool(workerCount int, bufferSize int, store storage.Storage, handlers *handler.Registry) *Pool {
	return &Pool{
		workerCount: workerCount,
		jobs:        make(chan *job.Job, bufferSize),
		storage:     store,
		handlers:    handlers,
	}
}

// Enqueue adds a job to the worker pool
func (p *Pool) Enqueue(j *job.Job) {
	p.jobs <- j
}

// Start begins processing jobs with workers
func (p *Pool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 1; i <= p.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			p.worker(ctx, workerID)
		}(i)
	}

	log.Printf("[POOL] Started %d workers", p.workerCount)
	wg.Wait()
	log.Println("[POOL] All workers stopped")
}

// worker processes jobs from the queue
func (p *Pool) worker(ctx context.Context, id int) {
	log.Printf("[WORKER-%d] Started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[WORKER-%d] Shutting down", id)
			return
		case j := <-p.jobs:
			p.processJob(ctx, id, j)
		}
	}
}

// processJob executes a single job
func (p *Pool) processJob(ctx context.Context, workerID int, j *job.Job) {
	log.Printf("[WORKER-%d] Processing job %s (type: %s, attempt: %d)",
		workerID, j.ID, j.Type, j.Attempts+1)

	// Mark job as running
	j.MarkRunning()
	if err := p.storage.Update(ctx, j); err != nil {
		log.Printf("[WORKER-%d] Failed to update job %s: %v", workerID, j.ID, err)
		return
	}

	// Get handler
	handlerFunc, err := p.handlers.Get(j.Type)
	if err != nil {
		j.MarkFailed(fmt.Errorf("no handler registered: %w", err))
		p.storage.Update(ctx, j)
		log.Printf("[WORKER-%d] Job %s FAILED: %v", workerID, j.ID, err)
		return
	}

	// Execute handler
	err = handlerFunc(ctx, j.Payload)

	if err != nil {
		j.MarkFailed(err)

		if j.CanRetry() {
			log.Printf("[WORKER-%d] Job %s failed (attempt %d/%d): %v - RETRYING",
				workerID, j.ID, j.Attempts, j.MaxRetry, err)

			// Re-queue with backoff
			go func(job *job.Job) {
				backoff := job.BackoffDuration()
				time.Sleep(backoff)
				p.jobs <- job
			}(j)
		} else {
			log.Printf("[WORKER-%d] Job %s FAILED permanently after %d attempts: %v",
				workerID, j.ID, j.Attempts, err)
		}
	} else {
		j.MarkCompleted()
		log.Printf("[WORKER-%d] Job %s COMPLETED successfully", workerID, j.ID)
	}

	// Save final state
	if err := p.storage.Update(ctx, j); err != nil {
		log.Printf("[WORKER-%d] Failed to update job %s: %v", workerID, j.ID, err)
	}
}
