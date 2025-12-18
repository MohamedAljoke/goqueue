package goqueue_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MohamedAljoke/goqueue"
	"github.com/MohamedAljoke/goqueue/internal/entity"
	"github.com/MohamedAljoke/goqueue/internal/storage"
	"github.com/MohamedAljoke/goqueue/internal/usecase"
)

func TestJobProcessing(t *testing.T) {

	t.Run("should create job with pending status", func(t *testing.T) {
		jobStorage := storage.NewMemoryStorage()
		jobStorage.ClearStorage()

		job := entity.NewJob(3)
		jobStorage.SaveJob(context.Background(), job)

		if job.Status != entity.StatusPending {
			t.Errorf("expected job status to be 'pending', got %s", job.Status)
		}
		retrieved, exists := jobStorage.GetJob(context.Background(), job.ID)
		if !exists {
			t.Errorf("expected job to exist in storage")
		}
		if retrieved.ID != job.ID {
			t.Errorf("expected job ID %s, got %s", job.ID, retrieved.ID)
		}
	})

	t.Run("should have status completed after processing", func(t *testing.T) {
		jobStorage := storage.NewMemoryStorage()
		jobStorage.ClearStorage()
		processor := usecase.NewJobProcessor(jobStorage)

		job := entity.NewJob(3)
		handler := func(ctx context.Context, payload map[string]interface{}) error {
			return nil
		}

		err := processor.ProcessJob(context.Background(), job, handler)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if job.Status != entity.StatusCompleted {
			t.Errorf("expected job status to be 'completed', got %s", job.Status)
		}
	})

	t.Run("should have status failed after processing with error", func(t *testing.T) {
		jobStorage := storage.NewMemoryStorage()
		jobStorage.ClearStorage()
		processor := usecase.NewJobProcessor(jobStorage)

		job := entity.NewJob(1)
		handler := func(ctx context.Context, payload map[string]interface{}) error {
			return fmt.Errorf("handler error")
		}

		err := processor.ProcessJob(context.Background(), job, handler)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if job.Status != entity.StatusFailed {
			t.Errorf("expected job status to be 'failed', got %s", job.Status)
		}
	})

	t.Run("should process job with registered handler", func(t *testing.T) {
		queue := goqueue.NewQueue()

		queue.RegisterHandler("email", func(ctx context.Context, payload map[string]interface{}) error {
			return nil
		})

		job, err := queue.SubmitJob(context.Background(), "email", map[string]interface{}{"to": "test@test.com"}, 3)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if job.Status != entity.StatusCompleted {
			t.Errorf("expected completed status, got %s", job.Status)
		}
	})

	t.Run("should return error for unregistered handler", func(t *testing.T) {
		queue := goqueue.NewQueue()

		job, err := queue.SubmitJob(context.Background(), "unknown", map[string]interface{}{}, 3)

		if err == nil {
			t.Error("expected error for unregistered handler")
		}
		if job != nil {
			t.Error("expected nil job")
		}
	})
}
