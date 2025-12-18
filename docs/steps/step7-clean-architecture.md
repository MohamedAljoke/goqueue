# Step 7: Clean Architecture Refactoring

## Problem

Our code had some architectural issues:
- Job entity was coupled to Storage (called `SaveJob()` directly)
- No clear separation between domain logic and infrastructure
- Package structure didn't reflect the layers
- Hard to test domain logic in isolation

## Solution

Refactor to **Clean Architecture** with clear layers:
- **Entity Layer**: Pure domain objects (Job, Status)
- **Use Case Layer**: Business logic orchestration (JobProcessor)
- **Infrastructure Layer**: External concerns (Storage)

## Key Changes

### 1. Package Restructure

```
internal/
├── entity/          # Domain layer
│   ├── job.go
│   └── job_test.go
├── storage/         # Infrastructure layer
│   └── job_repository.go
└── usecase/         # Use case layer
    └── job_processor.go
```

### 2. Pure Domain Entity

**internal/entity/job.go**
```go
package entity

import (
	"fmt"
	"slices"
	"time"
)

type (
	Status string
	Job    struct {
		ID      string
		Type    string
		Payload map[string]interface{}
		Status  Status
	}
)

const (
	StatusPending    = Status("pending")
	StatusProcessing = Status("processing")
	StatusCompleted  = Status("completed")
	StatusFailed     = Status("failed")
)

type HandlerFunc func(payload map[string]interface{}) error

func NewJob() *Job {
	job := &Job{
		ID:      fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Type:    "default",
		Payload: make(map[string]interface{}),
		Status:  StatusPending,  // ✅ No SaveJob() call
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
	job.Status = status  // ✅ No SaveJob() call
	return nil
}

func (job *Job) canTransition(from, to Status) bool {
	validTransitions := map[Status][]Status{
		StatusPending:    {StatusProcessing},
		StatusProcessing: {StatusCompleted, StatusFailed},
		StatusFailed:     {},
		StatusCompleted:  {},
	}
	return slices.Contains(validTransitions[from], to)
}
```

**Why it's better:**
- Job has ZERO dependencies ✅
- Pure domain logic ✅
- Easy to test in isolation ✅
- Single Responsibility Principle ✅

### 3. Storage Interface

**internal/storage/job_repository.go**
```go
package storage

import "github.com/MohamedAljoke/goqueue/internal/entity"

var jobStorage = make(map[string]*entity.Job)

type JobStorage interface {
	SaveJob(job *entity.Job)
	GetJob(id string) (*entity.Job, bool)
	ClearStorage()
}

type MemoryStorage struct{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (m *MemoryStorage) SaveJob(job *entity.Job) {
	jobStorage[job.ID] = job
}

func (m *MemoryStorage) GetJob(id string) (*entity.Job, bool) {
	job, exists := jobStorage[id]
	return job, exists
}

func (m *MemoryStorage) ClearStorage() {
	jobStorage = make(map[string]*entity.Job)
}
```

**Why interface:**
- Easy to swap implementations (memory → database) ✅
- Testable with mocks ✅
- Dependency Inversion Principle ✅

### 4. Use Case Orchestrator

**internal/usecase/job_processor.go**
```go
package usecase

import (
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

func (jp *JobProcessor) ProcessJob(job *entity.Job, handler entity.HandlerFunc) error {
	job.ChangeStatus(entity.StatusProcessing)
	jp.storage.SaveJob(job)  // ✅ Use case coordinates persistence

	if err := handler(job.Payload); err != nil {
		job.ChangeStatus(entity.StatusFailed)
		jp.storage.SaveJob(job)
		return fmt.Errorf("error handling process: %w", err)
	}

	job.ChangeStatus(entity.StatusCompleted)
	jp.storage.SaveJob(job)
	return nil
}
```

**Why use case:**
- One business operation = one use case ✅
- Orchestrates domain + infrastructure ✅
- Keeps entity pure ✅

### 5. Updated Tests

**goqueue_test.go**
```go
package goqueue_test

import (
	"fmt"
	"testing"

	"github.com/MohamedAljoke/goqueue/internal/entity"
	"github.com/MohamedAljoke/goqueue/internal/storage"
	"github.com/MohamedAljoke/goqueue/internal/usecase"
)

func TestJobProcessing(t *testing.T) {

	t.Run("should create job with pending status", func(t *testing.T) {
		jobStorage := storage.NewMemoryStorage()
		jobStorage.ClearStorage()

		job := entity.NewJob()
		jobStorage.SaveJob(job)  // ✅ Caller handles persistence

		if job.Status != entity.StatusPending {
			t.Errorf("expected job status to be 'pending', got %s", job.Status)
		}
		retrieved, exists := jobStorage.GetJob(job.ID)
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
		processor := usecase.NewJobProcessor(jobStorage)  // ✅ Inject dependency

		job := entity.NewJob()
		handler := func(payload map[string]interface{}) error {
			return nil
		}

		err := processor.ProcessJob(job, handler)
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

		job := entity.NewJob()
		handler := func(payload map[string]interface{}) error {
			return fmt.Errorf("handler error")
		}

		err := processor.ProcessJob(job, handler)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if job.Status != entity.StatusFailed {
			t.Errorf("expected job status to be 'failed', got %s", job.Status)
		}
	})
}
```

## What We Learned

### Clean Architecture Layers
- **Domain (Entity)**: Business rules, no dependencies
- **Use Case**: Application logic, orchestrates entities + infrastructure
- **Infrastructure**: External concerns (DB, APIs, etc.)

### Dependency Direction
```
Infrastructure → Use Case → Entity
   (depends on)  (depends on)  (depends on nothing)
```

### Separation of Concerns
- Entity = What is a Job?
- Storage = Where do we keep Jobs?
- Use Case = How do we process Jobs?

### Dependency Injection
- Pass dependencies through constructors
- Code to interfaces, not implementations
- Easy to test and swap implementations

## Benefits

✅ **Testability**: Each layer can be tested independently
✅ **Maintainability**: Clear responsibilities
✅ **Flexibility**: Easy to swap storage implementations
✅ **Scalability**: Ready to add more use cases

## What's Next

We're now ready to add:
- Handler registry (map job types → handlers)
- Background processing with channels
- Worker pool for concurrency
- Queue orchestrator to tie it all together

The clean architecture foundation makes adding these features much easier!
