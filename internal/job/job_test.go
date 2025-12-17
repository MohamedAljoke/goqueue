package job

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	j := New("test_job", map[string]interface{}{"key": "value"}, 3)

	if j.ID == "" {
		t.Error("expected job ID to be set")
	}

	if j.Type != "test_job" {
		t.Errorf("expected type 'test_job', got '%s'", j.Type)
	}

	if j.Status != StatusPending {
		t.Errorf("expected status pending, got %s", j.Status)
	}

	if j.Attempts != 0 {
		t.Errorf("expected 0 attempts, got %d", j.Attempts)
	}

	if j.MaxRetry != 3 {
		t.Errorf("expected max retry 3, got %d", j.MaxRetry)
	}
}

func TestMarkRunning(t *testing.T) {
	j := New("test", nil, 3)

	j.MarkRunning()

	if j.Status != StatusRunning {
		t.Errorf("expected status running, got %s", j.Status)
	}

	if j.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", j.Attempts)
	}
}

func TestMarkCompleted(t *testing.T) {
	j := New("test", nil, 3)
	j.MarkRunning()

	j.MarkCompleted()

	if j.Status != StatusCompleted {
		t.Errorf("expected status completed, got %s", j.Status)
	}

	if j.Error != "" {
		t.Errorf("expected error to be cleared, got '%s'", j.Error)
	}
}

func TestMarkFailed_WithRetry(t *testing.T) {
	j := New("test", nil, 3)
	j.MarkRunning()

	j.MarkFailed(errors.New("test error"))

	if j.Status != StatusPending {
		t.Errorf("expected status pending (for retry), got %s", j.Status)
	}

	if j.Error != "test error" {
		t.Errorf("expected error message, got '%s'", j.Error)
	}
}

func TestMarkFailed_NoRetry(t *testing.T) {
	j := New("test", nil, 1)
	j.MarkRunning() // attempt 1

	j.MarkFailed(errors.New("test error"))

	if j.Status != StatusFailed {
		t.Errorf("expected status failed, got %s", j.Status)
	}
}

func TestCanRetry(t *testing.T) {
	j := New("test", nil, 3)

	if !j.CanRetry() {
		t.Error("expected job to be retryable before max attempts")
	}

	j.Attempts = 3

	if j.CanRetry() {
		t.Error("expected job not to be retryable after max attempts")
	}
}

func TestStatusIsTerminal(t *testing.T) {
	tests := []struct {
		status   Status
		terminal bool
	}{
		{StatusPending, false},
		{StatusRunning, false},
		{StatusCompleted, true},
		{StatusFailed, true},
		{StatusCanceled, true},
	}

	for _, tt := range tests {
		if tt.status.IsTerminal() != tt.terminal {
			t.Errorf("status %s: expected terminal=%v, got %v",
				tt.status, tt.terminal, tt.status.IsTerminal())
		}
	}
}
