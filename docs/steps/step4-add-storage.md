# Step 4: Storage Layer

## Problem

Jobs only existed in memory during processing - no persistence or retrieval.

## Solution

Add simple in-memory storage using a map.

## Concepts Learned

- Maps in Go
- Package-level variables
- Data persistence patterns
- Unique ID generation with timestamps

## Implementation

### Add Package-Level Storage

```go
import (
    "fmt"
    "time"
)

var jobStorage = make(map[string]*Job)
```

Package-level variable shared across all code in the package.

### Generate Unique IDs

```go
func NewJob() *Job {
    job := &Job{
        ID:      fmt.Sprintf("job_%d", time.Now().UnixNano()),  // Unique!
        Type:    "default",
        Payload: make(map[string]interface{}),
        Status:  "pending",
    }
    SaveJob(job)  // Auto-save
    return job
}
```

### Storage Functions

```go
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

### Update Process to Persist

```go
func (job *Job) Process(handler HandlerFunc) error {
    job.Status = "processing"
    SaveJob(job)  // Persist state change

    err := handler(job.Payload)
    if err != nil {
        job.Status = "failed"
        SaveJob(job)  // Persist failure
        return fmt.Errorf("error handling process %w", err)
    }

    job.Status = "completed"
    SaveJob(job)  // Persist completion
    return nil
}
```

## Tests

### Test Storage Persistence

```go
t.Run("should persist status changes after processing", func(t *testing.T) {
    goqueue.ClearStorage()  // Clean slate

    job := goqueue.NewJob()
    handler := func(payload map[string]interface{}) error {
        return nil
    }

    job.Process(handler)

    // Retrieve from storage
    retrieved, exists := goqueue.GetJob(job.ID)
    if !exists {
        t.Errorf("expected job to exist")
    }
    if retrieved.Status != "completed" {
        t.Errorf("expected status 'completed' in storage, got %s", retrieved.Status)
    }
})
```

### Test Retrieval

```go
t.Run("should save and retrieve job", func(t *testing.T) {
    goqueue.ClearStorage()

    job := goqueue.NewJob()

    retrieved, exists := goqueue.GetJob(job.ID)
    if !exists {
        t.Errorf("expected to find job")
    }
    if retrieved.ID != job.ID {
        t.Errorf("expected job ID %s, got %s", job.ID, retrieved.ID)
    }
})
```

## Key Concepts

### Maps in Go

```go
// Create
m := make(map[string]*Job)

// Set
m["job_123"] = job

// Get
job, exists := m["job_123"]
if exists {
    // job found
}

// Delete
delete(m, "job_123")
```

### Multiple Return Values

```go
func GetJob(id string) (*Job, bool) {
    job, exists := jobStorage[id]
    return job, exists
}

// Usage
job, exists := GetJob("job_123")
if exists {
    fmt.Println(job.Status)
}
```

Common Go pattern:
- First return: the value
- Second return: whether it was found

### Unique ID Generation

```go
fmt.Sprintf("job_%d", time.Now().UnixNano())
```

- `time.Now().UnixNano()` returns nanoseconds since epoch
- Highly unlikely to collide
- Sortable by creation time

Example IDs:
```
job_1702850123456789012
job_1702850123456789345
job_1702850123456789678
```

### Package-Level Variables

```go
var jobStorage = make(map[string]*Job)
```

- Shared across entire package
- Initialized once when package loads
- Accessible by all functions in package

## Testing with Shared State

```go
func TestStorage(t *testing.T) {
    t.Run("test 1", func(t *testing.T) {
        goqueue.ClearStorage()  // Clean state
        // Test...
    })

    t.Run("test 2", func(t *testing.T) {
        goqueue.ClearStorage()  // Clean state
        // Test...
    })
}
```

Always clear shared state between tests!

## In-Memory Storage Limitations

### Pros
✅ Simple to implement
✅ Fast
✅ Good for development/testing

### Cons
❌ Lost on restart
❌ Not shared across processes
❌ Limited by memory

**Solution:** Use interface pattern to swap storage backends later (Redis, PostgreSQL, etc.)

## What's Still Wrong?

1. **Tight coupling** - Job directly calls SaveJob()
2. **No separation** - Storage and job logic mixed
3. **Single file** - Everything in one place

## Next Step

[Step 5: Package Organization](./step5-refactor-internal-package.md) - Split into logical packages
