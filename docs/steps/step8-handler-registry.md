# Step 8: Handler Registry

## Problem

In the previous steps, we had to pass the handler function every time we processed a job:
```go
processor.ProcessJob(job, emailHandler)
processor.ProcessJob(job, smsHandler)
```

This is repetitive and doesn't scale well. In a real job queue:
- Different job types need different handlers
- Handlers should be registered once, used many times
- The system should automatically route jobs to the correct handler

## Solution

Create a **Handler Registry** that maps job types to their handlers:
```go
queue.RegisterHandler("email", emailHandler)  // Register once
queue.RegisterHandler("sms", smsHandler)

queue.SubmitJob("email", payload)  // Automatically uses emailHandler
queue.SubmitJob("sms", payload)    // Automatically uses smsHandler
```

## Implementation

### 1. Create Handler Registry

**internal/handler/registry.go**
```go
package handler

import (
	"fmt"

	"github.com/MohamedAljoke/goqueue/internal/entity"
)

type Registry struct {
	handlers map[string]entity.HandlerFunc
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]entity.HandlerFunc),
	}
}

func (r *Registry) RegisterHandler(jobType string, handler entity.HandlerFunc) {
	r.handlers[jobType] = handler
}

func (r *Registry) GetHandler(jobType string) (entity.HandlerFunc, error) {
	handler, exists := r.handlers[jobType]
	if !exists {
		return nil, fmt.Errorf("no handler registered for job type: %s", jobType)
	}
	return handler, nil
}
```

**Key concepts:**
- `map[string]entity.HandlerFunc` - Maps job type → handler function
- `RegisterHandler()` - Stores a handler for a job type
- `GetHandler()` - Retrieves handler, returns error if not found

### 2. Integrate Registry into Queue

**goqueue.go**
```go
package goqueue

import (
	"github.com/MohamedAljoke/goqueue/internal/entity"
	"github.com/MohamedAljoke/goqueue/internal/handler"
	"github.com/MohamedAljoke/goqueue/internal/storage"
	"github.com/MohamedAljoke/goqueue/internal/usecase"
)

type Queue struct {
	storage   storage.JobStorage
	processor *usecase.JobProcessor
	registry  *handler.Registry  // ← Added
}

func NewQueue() *Queue {
	storage := storage.NewMemoryStorage()
	processor := usecase.NewJobProcessor(storage)
	registry := handler.NewRegistry()  // ← Create registry

	return &Queue{
		storage:   storage,
		processor: processor,
		registry:  registry,  // ← Add to struct
	}
}

// Public API: Register a handler
func (q *Queue) RegisterHandler(jobType string, handlerFunc entity.HandlerFunc) {
	q.registry.RegisterHandler(jobType, handlerFunc)
}

// Updated: No handler parameter, looks it up from registry
func (q *Queue) SubmitJob(jobType string, payload map[string]interface{}) (*entity.Job, error) {
	// Look up handler from registry
	handler, err := q.registry.GetHandler(jobType)
	if err != nil {
		return nil, err  // No handler registered for this job type
	}

	job := entity.NewJob()
	job.Type = jobType
	job.Payload = payload

	return job, q.processor.ProcessJob(job, handler)
}
```

**Changes:**
- Added `registry` field to Queue
- New `RegisterHandler()` method (public API)
- `SubmitJob()` no longer takes handler parameter
- `SubmitJob()` looks up handler from registry automatically

### 3. Updated Tests

**goqueue_test.go**
```go
t.Run("should process job with registered handler", func(t *testing.T) {
	queue := goqueue.NewQueue()

	// Register handler first
	queue.RegisterHandler("email", func(payload map[string]interface{}) error {
		return nil
	})

	// Submit job (no handler parameter)
	job, err := queue.SubmitJob("email", map[string]interface{}{"to": "test@test.com"})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if job.Status != entity.StatusCompleted {
		t.Errorf("expected completed status, got %s", job.Status)
	}
})

t.Run("should return error for unregistered handler", func(t *testing.T) {
	queue := goqueue.NewQueue()

	// Try to submit without registering handler
	job, err := queue.SubmitJob("unknown", map[string]interface{}{})

	if err == nil {
		t.Error("expected error for unregistered handler")
	}
	if job != nil {
		t.Error("expected nil job")
	}
})
```

## What We Learned

### Maps in Go
```go
myMap := make(map[string]HandlerFunc)  // Create
myMap["key"] = value                   // Store
value, exists := myMap["key"]          // Retrieve with existence check
```

### Error Handling Pattern
```go
value, exists := myMap[key]
if !exists {
    return nil, fmt.Errorf("not found: %s", key)
}
return value, nil
```

### Registry Pattern
- **Single source of truth** for handler mappings
- **Register once, use many times**
- **Automatic routing** based on job type
- **Type safety** through Go's type system

### Separation of Concerns
- **Registry**: Knows about job type → handler mapping
- **Queue**: Orchestrates registration + job submission
- **Use Case**: Processes jobs with handlers
- **Entity**: Pure domain logic

## Benefits

✅ **DRY (Don't Repeat Yourself)**: Register handlers once
✅ **Clean API**: `SubmitJob()` is simpler, no handler parameter
✅ **Type Safety**: Error if handler not registered
✅ **Scalability**: Easy to add new job types
✅ **Testability**: Can test handler registration independently

## Note on Thread Safety

**Current implementation:** Not thread-safe (no mutex)

**Why it's OK for now:** We're still processing jobs synchronously in a single goroutine.

**When we'll need mutex:** In Step 9-10 when we add:
- Background workers (multiple goroutines)
- Concurrent access to the registry
- Race conditions between reading/writing the map

We'll add `sync.RWMutex` then!

## What's Next

Now that we have handler registry, we're ready for the exciting part:
- **Step 9**: Background processing with channels
- **Step 10**: Worker pool for concurrency

The registry makes async processing much cleaner - workers can look up handlers automatically!
