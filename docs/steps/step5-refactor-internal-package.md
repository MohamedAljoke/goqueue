# Step 5: Package Organization

## Problem

All code was in one file - needed better organization and separation of concerns.

## Solution

Split into logical packages:
- `internal/job.go` - Job domain logic
- `internal/storage.go` - Storage operations
- `goqueue.go` - Public API

## Concepts Learned

- Package organization in Go
- Internal package convention
- Separation of concerns
- Interface definitions

## File Structure

```
goqueue/
├── goqueue.go              # Public API
├── goqueue_test.go         # Tests
├── internal/
│   ├── job.go             # Job logic
│   └── storage.go         # Storage logic
└── go.mod
```

## Implementation

### internal/job.go

```go
package internal

import (
    "fmt"
    "time"
)

type Job struct {
    ID      string
    Type    string
    Payload map[string]interface{}
    Status  string
}

type HandlerFunc func(payload map[string]interface{}) error

func NewJob() *Job {
    job := &Job{
        ID:      fmt.Sprintf("job_%d", time.Now().UnixNano()),
        Type:    "default",
        Payload: make(map[string]interface{}),
        Status:  "pending",
    }
    SaveJob(job)
    return job
}

func (job *Job) Process(handler HandlerFunc) error {
    job.Status = "processing"
    SaveJob(job)

    err := handler(job.Payload)
    if err != nil {
        job.Status = "failed"
        SaveJob(job)
        return fmt.Errorf("error handling process %w", err)
    }

    job.Status = "completed"
    SaveJob(job)
    return nil
}
```

### internal/storage.go

```go
package internal

var jobStorage = make(map[string]*Job)

// Storage interface for future implementations
type Storage interface {
    SaveJob(job *Job)
    GetJob(id string) (*Job, bool)
    ClearStorage()
}

type MemoryStorage struct{}

// Package-level functions (not using interface yet)
func SaveJob(job *Job) {
    jobStorage[job.ID] = job
}

func GetJob(id string) (*Job, bool) {
    job, exists := jobStorage[id]
    return job, exists
}

func ClearStorage() {
    jobStorage = make(map[string]*Job)
}
```

### goqueue.go (Public API)

```go
package goqueue

import (
    "github.com/MohamedAljoke/goqueue/internal"
)

func NewJob() *internal.Job {
    return internal.NewJob()
}
```

### Updated Tests

```go
package goqueue_test

import (
    "testing"
    "github.com/MohamedAljoke/goqueue"
    "github.com/MohamedAljoke/goqueue/internal"
)

func TestGenerateJob(t *testing.T) {
    t.Run("should create job with pending status", func(t *testing.T) {
        internal.ClearStorage()

        job := goqueue.NewJob()
        if job.Status != "pending" {
            t.Errorf("expected job status to be 'pending', got %s", job.Status)
        }

        retrieved, exists := internal.GetJob(job.ID)
        if !exists {
            t.Errorf("expected job to exist in storage")
        }
        if retrieved.ID != job.ID {
            t.Errorf("expected job ID %s, got %s", job.ID, retrieved.ID)
        }
    })
}
```

## Key Concepts

### Package Organization

```
package internal  // Package declaration

import (...)      // Import other packages

type Job struct{} // Type definitions

func NewJob()     // Functions
```

### Internal Package Convention

`internal/` has special meaning in Go:
- Can only be imported by parent package and siblings
- Cannot be imported by external packages
- Enforces encapsulation

```go
// OK: Same module
import "github.com/MohamedAljoke/goqueue/internal"

// ERROR: External module
import "github.com/someone-else/goqueue/internal"
```

### Separation of Concerns

Each package has one responsibility:

| Package | Responsibility |
|---------|---------------|
| `internal/job` | Job lifecycle and state |
| `internal/storage` | Data persistence |
| `goqueue` | Public API |

### Interface Definition

```go
type Storage interface {
    SaveJob(job *Job)
    GetJob(id string) (*Job, bool)
    ClearStorage()
}
```

Even though not used yet, defining interfaces early:
- Documents intended behavior
- Makes future refactoring easier
- Enables testing with mocks

## Import Paths

### Module Declaration (go.mod)

```
module github.com/MohamedAljoke/goqueue
```

### Importing

```go
// Import by full path
import "github.com/MohamedAljoke/goqueue/internal"

// Use package name
job := internal.NewJob()
```

## Benefits of This Structure

### Before (Single File)
❌ Everything mixed together
❌ Hard to find code
❌ Difficult to test in isolation
❌ No clear boundaries

### After (Multiple Packages)
✅ Clear separation
✅ Easy to navigate
✅ Can test each package independently
✅ Logical boundaries

## What's Still Wrong?

1. **Storage interface not used** - Defined but functions are package-level
2. **Job coupled to storage** - Calls SaveJob() directly
3. **String-based status** - No type safety
4. **Tests import internal** - Should use public API

## Next Step

[Step 6: State Machine Pattern](./step6-state-machine.md) - Type-safe status management
