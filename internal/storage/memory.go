package storage

import (
	"context"
	"sync"

	"github.com/MohamedAljoke/goqueue/internal/job"
)

// Memory is an in-memory implementation of Storage
type Memory struct {
	mu   sync.RWMutex
	jobs map[string]*job.Job
}

// NewMemory creates a new in-memory storage
func NewMemory() *Memory {
	return &Memory{
		jobs: make(map[string]*job.Job),
	}
}

// Save persists a job in memory
func (m *Memory) Save(ctx context.Context, j *job.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobs[j.ID] = j
	return nil
}

// Get retrieves a job by ID
func (m *Memory) Get(ctx context.Context, id string) (*job.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	j, ok := m.jobs[id]
	if !ok {
		return nil, ErrNotFound
	}

	return j, nil
}

// List retrieves jobs by status
func (m *Memory) List(ctx context.Context, status job.Status) ([]*job.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*job.Job
	for _, j := range m.jobs {
		if j.Status == status {
			result = append(result, j)
		}
	}

	return result, nil
}

// Update updates an existing job
func (m *Memory) Update(ctx context.Context, j *job.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.jobs[j.ID]; !ok {
		return ErrNotFound
	}

	m.jobs[j.ID] = j
	return nil
}

// Delete removes a job
func (m *Memory) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.jobs[id]; !ok {
		return ErrNotFound
	}

	delete(m.jobs, id)
	return nil
}
