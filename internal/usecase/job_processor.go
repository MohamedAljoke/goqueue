package usecase

import (
	"context"
	"fmt"

	"github.com/MohamedAljoke/goqueue/internal/entity"
)

type JobProcessor struct {
}

func NewJobProcessor() *JobProcessor {
	return &JobProcessor{}
}

func (jp *JobProcessor) Process(tx context.Context, job *entity.Job, handler entity.HandlerFunc) error {
	if err := job.MarkRunning(); err != nil {
		return err
	}
	if err := handler(tx, job.Payload); err != nil {
		if markErr := job.MarkFailed(err); markErr != nil {
			return markErr
		}
		return fmt.Errorf("error handling process: %w", err)
	}
	if err := job.MarkCompleted(); err != nil {
		return err
	}

	return nil
}
