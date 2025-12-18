package storage

import (
	"context"
	"fmt"

	"github.com/MohamedAljoke/goqueue/internal/entity"
)

var jobStorage = make(map[string]*entity.Job)

type JobStorage interface {
	SaveJob(ctx context.Context, job *entity.Job)
	GetJob(ctx context.Context, id string) (*entity.Job, bool)
	UpdateJob(ctx context.Context, job *entity.Job) error
	ClearStorage()
}

type MemoryStorage struct{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (m *MemoryStorage) SaveJob(tx context.Context, job *entity.Job) {
	jobStorage[job.ID] = job
}
func (m *MemoryStorage) GetJob(tx context.Context, id string) (*entity.Job, bool) {
	job, exists := jobStorage[id]
	return job, exists
}
func (m *MemoryStorage) UpdateJob(ctx context.Context, job *entity.Job) error {
	if _, exists := jobStorage[job.ID]; !exists {
		return fmt.Errorf("job %s not found", job.ID)
	}
	jobStorage[job.ID] = job
	return nil
}
func (m *MemoryStorage) ClearStorage() {
	jobStorage = make(map[string]*entity.Job)
}
