package usecase

import (
	"context"
	"fmt"

	"github.com/MohamedAljoke/goqueue/internal/entity"
	"github.com/MohamedAljoke/goqueue/internal/storage"
)

type JobProcessor struct {
	storage storage.JobStorage
}

func NewJobProcessor(storage storage.JobStorage) *JobProcessor {
	return &JobProcessor{storage}
}

func (jp *JobProcessor) ProcessJob(tx context.Context, job *entity.Job, handler entity.HandlerFunc) error {
	if err := job.MarkRunning(); err != nil {
		return err
	}
	jp.storage.SaveJob(tx, job)

	if err := handler(tx, job.Payload); err != nil {
		if markErr := job.MarkFailed(err); markErr != nil {
			return markErr
		}
		jp.storage.SaveJob(tx, job)
		return fmt.Errorf("error handling process: %w", err)
	}

	if err := job.MarkCompleted(); err != nil {
		return err
	}
	jp.storage.SaveJob(tx, job)
	return nil
}
