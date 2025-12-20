package steps_test

import (
	"fmt"
	"testing"

	"github.com/MohamedAljoke/goqueue/steps"
)

func TestGeneric(t *testing.T) {
	t.Run("should have status completed after processing", func(t *testing.T) {
		job := steps.NewJob()
		handler := func() error {
			return nil
		}

		err := job.Process(handler)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if job.Status != "completed" {
			t.Errorf("expected job status to be 'completed', got %s", job.Status)
		}
	})
	t.Run("should have status failed after processing with error", func(t *testing.T) {
		job := steps.NewJob()
		handler := func() error {
			return fmt.Errorf("some error")
		}

		err := job.Process(handler)

		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if job.Status != "failed" {
			t.Errorf("expected job status to be 'failed', got %s", job.Status)
		}
	})
}
