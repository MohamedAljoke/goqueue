package entity

import (
	"context"
	"fmt"
	"slices"
	"time"
)

type (
	Status string
	Job    struct {
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
)

const (
	StatusPending    = Status("pending")
	StatusProcessing = Status("processing")
	StatusCompleted  = Status("completed")
	StatusFailed     = Status("failed")
)

type HandlerFunc func(ctx context.Context, payload map[string]interface{}) error

func NewJob(maxRetry int) *Job {
	now := time.Now()
	job := &Job{
		ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Status:    StatusPending,
		Attempts:  0,
		MaxRetry:  maxRetry,
		Error:     "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	return job
}

func (job *Job) ChangeStatus(status Status) error {
	if !job.canTransition(job.Status, status) {
		return fmt.Errorf(
			"invalid job status transition: %s -> %s",
			job.Status,
			status,
		)
	}

	job.Status = status

	return nil
}

func (job *Job) canTransition(from, to Status) bool {
	validTransitions := map[Status][]Status{
		StatusPending: {
			StatusProcessing,
		},
		StatusProcessing: {
			StatusCompleted,
			StatusFailed,
			StatusPending, // NEW: Allow retry (back to pending)
		},
		StatusFailed:    {},
		StatusCompleted: {},
	}

	return slices.Contains(validTransitions[from], to)
}

func (j *Job) MarkRunning() error {
	if err := j.ChangeStatus(StatusProcessing); err != nil {
		return err
	}
	j.Attempts++
	j.UpdatedAt = time.Now()
	return nil
}

func (j *Job) MarkCompleted() error {
	if err := j.ChangeStatus(StatusCompleted); err != nil {
		return err
	}
	j.Error = ""
	j.UpdatedAt = time.Now()
	return nil
}

func (j *Job) MarkFailed(err error) error {
	j.Error = err.Error()
	j.UpdatedAt = time.Now()

	var targetStatus Status
	if j.CanRetry() {
		targetStatus = StatusPending
	} else {
		targetStatus = StatusFailed
	}

	return j.ChangeStatus(targetStatus)
}

func (j *Job) CanRetry() bool {
	return j.Attempts < j.MaxRetry
}
