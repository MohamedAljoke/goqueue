package job_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/MohamedAljoke/goqueue/internal/job"
)

func TestJob_StatusTransition(t *testing.T) {

	t.Run("should create a new job with pending status", func(t *testing.T) {
		newJob := job.NewJob(3)

		if newJob.Status != job.StatusPending {
			t.Errorf("expected job status to be pending, got %s", newJob.Status)
		}
		if newJob.Attempts != 0 {
			t.Errorf("expected job attempts to be 0, got %d", newJob.Attempts)
		}
		if newJob.MaxRetry != 3 {
			t.Errorf("expected job max retry to be 3, got %d", newJob.MaxRetry)
		}
	})
	t.Run("should change job status to processing then completed correctly", func(t *testing.T) {
		newJob := job.NewJob(3)

		err := newJob.ChangeStatus(job.StatusProcessing)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if newJob.Status != job.StatusProcessing {
			t.Errorf("expected job status to be processing, got %s", newJob.Status)
		}

		err = newJob.ChangeStatus(job.StatusCompleted)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if newJob.Status != job.StatusCompleted {
			t.Errorf("expected job status to be completed, got %s", newJob.Status)
		}
	})
	t.Run("should not allow invalid status transitions", func(t *testing.T) {
		newJob := job.NewJob(3)

		err := newJob.ChangeStatus(job.StatusCompleted)
		if err == nil {
			t.Fatalf("expected error for invalid status transition, got nil")
		}

		if !errors.Is(err, job.ErrInvalidStatusTransition) {
			t.Fatalf("expected ErrInvalidStatusTransition, got %v", err)
		}

		err = newJob.ChangeStatus(job.StatusProcessing)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

}

func TestJob_RetryLogic(t *testing.T) {
	t.Run("should allow retry when attempts are less than max retry", func(t *testing.T) {
		newJob := job.NewJob(3)
		newJob.Attempts = 2

		if !newJob.CanRetry() {
			t.Fatalf("expected can retry to be true")
		}
	})
	t.Run("should not allow retry when attempts reach max retry", func(t *testing.T) {
		newJob := job.NewJob(3)
		newJob.Attempts = 3

		if newJob.CanRetry() {
			t.Fatalf("expected can retry to be false")
		}
	})
}

func TestJob_BackoffDuration(t *testing.T) {
	tests := []struct {
		attempts int
		expected int
	}{
		{attempts: 0, expected: 0},
		{attempts: 1, expected: 1},
		{attempts: 2, expected: 4},
		{attempts: 3, expected: 9},
		{attempts: 4, expected: 16},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempts_%d", tt.attempts), func(t *testing.T) {
			j := job.NewJob(3)
			j.Attempts = tt.attempts

			d := j.BackoffDuration()
			if d != time.Duration(tt.expected)*time.Second {
				t.Fatalf(
					"expected backoff %ds, got %s",
					tt.expected,
					d,
				)
			}
		})
	}
}

func TestJob_MarkRunning(t *testing.T) {
	t.Run("should move job to processing and increment attempts", func(t *testing.T) {
		j := job.NewJob(3)

		err := j.MarkRunning()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if j.Status != job.StatusProcessing {
			t.Fatalf("expected status processing, got %s", j.Status)
		}

		if j.Attempts != 1 {
			t.Fatalf("expected attempts to be 1, got %d", j.Attempts)
		}
	})
}

func TestJob_MarkCompleted(t *testing.T) {
	t.Run("should mark job as completed and clear error", func(t *testing.T) {
		j := job.NewJob(3)

		_ = j.MarkRunning()
		j.Error = "some error"

		err := j.MarkCompleted()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if j.Status != job.StatusCompleted {
			t.Fatalf("expected status completed, got %s", j.Status)
		}

		if j.Error != "" {
			t.Fatalf("expected error to be cleared")
		}
	})
}

func TestJob_MarkFailed_WithRetry(t *testing.T) {
	t.Run("should return to pending when retry is available", func(t *testing.T) {
		j := job.NewJob(3)

		_ = j.MarkRunning()

		err := j.MarkFailed(errors.New("boom"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if j.Status != job.StatusPending {
			t.Fatalf("expected status pending, got %s", j.Status)
		}

		if j.Error == "" {
			t.Fatalf("expected error to be set")
		}
	})
}

func TestJob_MarkFailed_NoRetry(t *testing.T) {
	t.Run("should move job to failed when retries are exhausted", func(t *testing.T) {
		j := job.NewJob(1)

		_ = j.MarkRunning()

		err := j.MarkFailed(errors.New("boom"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if j.Status != job.StatusFailed {
			t.Fatalf("expected status failed, got %s", j.Status)
		}
	})
}
