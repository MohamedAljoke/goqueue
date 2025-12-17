package job

import (
	"fmt"
	"time"
)

// Status represents the current state of a job
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCanceled  Status = "canceled"
)

// IsTerminal returns true if the status is a terminal state
func (s Status) IsTerminal() bool {
	return s == StatusCompleted || s == StatusFailed || s == StatusCanceled
}

// Job represents a unit of work to be processed
type Job struct {
	ID        string
	Type      string
	Payload   map[string]interface{}
	Status    Status
	Attempts  int
	MaxRetry  int
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New creates a new job
func New(jobType string, payload map[string]interface{}, maxRetry int) *Job {
	now := time.Now()
	return &Job{
		ID:        generateID(),
		Type:      jobType,
		Payload:   payload,
		Status:    StatusPending,
		Attempts:  0,
		MaxRetry:  maxRetry,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MarkRunning transitions the job to running status
func (j *Job) MarkRunning() {
	j.Status = StatusRunning
	j.Attempts++
	j.UpdatedAt = time.Now()
}

// MarkCompleted transitions the job to completed status
func (j *Job) MarkCompleted() {
	j.Status = StatusCompleted
	j.Error = ""
	j.UpdatedAt = time.Now()
}

// MarkFailed transitions the job to failed status
func (j *Job) MarkFailed(err error) {
	j.Error = err.Error()
	j.UpdatedAt = time.Now()

	if j.CanRetry() {
		j.Status = StatusPending
	} else {
		j.Status = StatusFailed
	}
}

// MarkCanceled transitions the job to canceled status
func (j *Job) MarkCanceled() {
	j.Status = StatusCanceled
	j.UpdatedAt = time.Now()
}

// CanRetry returns true if the job can be retried
func (j *Job) CanRetry() bool {
	return j.Attempts < j.MaxRetry
}

// BackoffDuration returns the duration to wait before retrying
func (j *Job) BackoffDuration() time.Duration {
	// Exponential backoff: attempt² seconds
	backoff := j.Attempts * j.Attempts
	return time.Duration(backoff) * time.Second
}

// generateID creates a unique job ID
func generateID() string {
	return fmt.Sprintf("job_%d", time.Now().UnixNano())
}
