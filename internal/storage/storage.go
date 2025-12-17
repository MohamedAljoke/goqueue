package storage

import (
	"context"
	"errors"

	"github.com/MohamedAljoke/goqueue/internal/job"
)

var (
	ErrNotFound = errors.New("job not found")
)

// Storage defines the interface for job persistence
type Storage interface {
	// Save persists a job
	Save(ctx context.Context, j *job.Job) error

	// Get retrieves a job by ID
	Get(ctx context.Context, id string) (*job.Job, error)

	// List retrieves jobs by status
	List(ctx context.Context, status job.Status) ([]*job.Job, error)

	// Update updates an existing job
	Update(ctx context.Context, j *job.Job) error

	// Delete removes a job
	Delete(ctx context.Context, id string) error
}
