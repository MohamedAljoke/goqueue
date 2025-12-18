# Step 9: Context Support & Job Retry

## Problem

Our current implementation has limitations:
1. **No retry mechanism** - Jobs fail permanently on first error
2. **No context** - Can't cancel jobs or set timeouts
3. **No error tracking** - Don't know why jobs failed
4. **No timestamps** - Can't audit when jobs were created/updated

Real-world job queues need:
- Automatic retry with exponential backoff
- Graceful shutdown (context cancellation)
- Error messages for debugging
- Audit trail (timestamps)

## Solution

Add retry logic and context support to all layers.

## Implementation

### 1. Add Retry Fields to Job Entity

**internal/entity/job.go**
```go
type Job struct {
    ID        string
    Type      string
    Payload   map[string]interface{}
    Status    Status
    Attempts  int        // NEW: How many times we've tried
    MaxRetry  int        // NEW: Maximum retry attempts
    Error     string     // NEW: Last error message
    CreatedAt time.Time  // NEW: When job was created
    UpdatedAt time.Time  // NEW: Last update timestamp
}
```

**Update NewJob() to accept maxRetry:**
```go
func NewJob(maxRetry int) *Job {
    now := time.Now()
    return &Job{
        ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
        Status:    StatusPending,
        Attempts:  0,
        MaxRetry:  maxRetry,
        Error:     "",
        CreatedAt: now,
        UpdatedAt: now,
    }
}
```

### 2. Add Retry Methods to Job

**CanRetry() - Check if job can be retried:**
```go
func (j *Job) CanRetry() bool {
    return j.Attempts < j.MaxRetry
}
```

**BackoffDuration() - Exponential backoff calculation:**
```go
func (j *Job) BackoffDuration() time.Duration {
    // Exponential backoff: attempt² seconds
    backoff := j.Attempts * j.Attempts
    return time.Duration(backoff) * time.Second
}
```
- Attempt 1 fails → wait 1 second
- Attempt 2 fails → wait 4 seconds
- Attempt 3 fails → wait 9 seconds

**Update state transition rules to allow retry:**
```go
func (job *Job) canTransition(from, to Status) bool {
    validTransitions := map[Status][]Status{
        StatusPending: {
            StatusProcessing,
        },
        StatusProcessing: {
            StatusCompleted,
            StatusFailed,
            StatusPending,  // NEW: Allow retry (back to pending)
        },
        StatusFailed:    {},
        StatusCompleted: {},
    }
    return slices.Contains(validTransitions[from], to)
}
```

### 3. Add Mark Methods

Instead of manually calling `ChangeStatus()`, these methods bundle state transition + metadata updates:

**MarkRunning() - Start processing:**
```go
func (j *Job) MarkRunning() error {
    if err := j.ChangeStatus(StatusProcessing); err != nil {
        return err
    }
    j.Attempts++  // Increment attempt counter
    j.UpdatedAt = time.Now()
    return nil
}
```

**MarkCompleted() - Job succeeded:**
```go
func (j *Job) MarkCompleted() error {
    if err := j.ChangeStatus(StatusCompleted); err != nil {
        return err
    }
    j.Error = ""  // Clear any previous error
    j.UpdatedAt = time.Now()
    return nil
}
```

**MarkFailed() - Job failed (retry or terminate):**
```go
func (j *Job) MarkFailed(err error) error {
    j.Error = err.Error()
    j.UpdatedAt = time.Now()

    var targetStatus Status
    if j.CanRetry() {
        targetStatus = StatusPending  // Retry: back to pending
    } else {
        targetStatus = StatusFailed   // Give up: terminal failure
    }

    return j.ChangeStatus(targetStatus)
}
```

**Key insight:** `MarkFailed()` uses `CanRetry()` to decide:
- If `Attempts < MaxRetry` → set status to `Pending` (will retry)
- Otherwise → set status to `Failed` (terminal state)

### 4. Update HandlerFunc Signature

**Add context.Context parameter:**
```go
type HandlerFunc func(ctx context.Context, payload map[string]interface{}) error
```

**Why context:**
- Cancellation: Stop long-running jobs
- Timeouts: Job must complete in X seconds
- Request-scoped values
- Essential for worker graceful shutdown

### 5. Update Storage Layer

**Add context to all storage methods:**

```go
type JobStorage interface {
    SaveJob(ctx context.Context, job *entity.Job)
    GetJob(ctx context.Context, id string) (*entity.Job, bool)
    ClearStorage()
}

func (m *MemoryStorage) SaveJob(ctx context.Context, job *entity.Job) {
    jobStorage[job.ID] = job
}

func (m *MemoryStorage) GetJob(ctx context.Context, id string) (*entity.Job, bool) {
    job, exists := jobStorage[id]
    return job, exists
}
```

For now we just accept `ctx` parameter - we'll use it when we add workers.

### 6. Update JobProcessor Use Case

**Update ProcessJob to use Mark methods:**

```go
func (jp *JobProcessor) ProcessJob(ctx context.Context, job *entity.Job, handler entity.HandlerFunc) error {
    // Mark job as running (increments Attempts)
    if err := job.MarkRunning(); err != nil {
        return err
    }
    jp.storage.SaveJob(ctx, job)

    // Execute handler with context
    if err := handler(ctx, job.Payload); err != nil {
        // Mark as failed (will retry or terminate based on MaxRetry)
        if markErr := job.MarkFailed(err); markErr != nil {
            return markErr
        }
        jp.storage.SaveJob(ctx, job)
        return fmt.Errorf("handler failed: %w", err)
    }

    // Mark as completed
    if err := job.MarkCompleted(); err != nil {
        return err
    }
    jp.storage.SaveJob(ctx, job)
    return nil
}
```

**Benefits:**
- `MarkRunning()` automatically increments `Attempts`
- `MarkFailed()` automatically decides retry vs terminal failure
- All state changes go through state machine validation
- Storage saves updated timestamps and error messages

### 7. Update Queue Public API

**Add context and maxRetry to SubmitJob:**

```go
func (q *Queue) SubmitJob(ctx context.Context, jobType string, payload map[string]interface{}, maxRetry int) (*entity.Job, error) {
    handler, err := q.registry.GetHandler(jobType)
    if err != nil {
        return nil, err
    }

    job := entity.NewJob(maxRetry)  // Pass maxRetry
    job.Type = jobType
    job.Payload = payload

    return job, q.processor.ProcessJob(ctx, job, handler)
}
```

### 8. Update Tests

**All tests need:**
- Import `context`
- Pass `context.Background()` to all methods
- Pass `maxRetry` to `NewJob()` and `SubmitJob()`
- Update handler signatures to include `ctx`

**Example:**
```go
t.Run("should have status failed after processing with error", func(t *testing.T) {
    jobStorage := storage.NewMemoryStorage()
    jobStorage.ClearStorage()
    processor := usecase.NewJobProcessor(jobStorage)

    job := entity.NewJob(1)  // maxRetry=1, will fail immediately
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
```

## What We Learned

### Context Pattern in Go
```go
func DoWork(ctx context.Context, params ...) error {
    // Check if context was cancelled
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Do work...
    }
}
```

Context is first parameter by convention.

### Exponential Backoff
```go
backoff := attempts * attempts  // 1, 4, 9, 16, 25...
time.Sleep(time.Duration(backoff) * time.Second)
```

Prevents overwhelming failing services.

### Retry State Machine
```
Pending → Processing → Failed (if CanRetry()) → Pending → Processing → ...
                     → Failed (terminal, no more retries)
                     → Completed
```

### Mark Methods vs Manual State Changes

**Before:**
```go
job.ChangeStatus(StatusProcessing)
job.Attempts++
job.UpdatedAt = time.Now()
storage.SaveJob(job)
```

**After:**
```go
job.MarkRunning()  // Handles everything!
storage.SaveJob(ctx, job)
```

### Why Timestamps Matter

- `CreatedAt` - When job entered system
- `UpdatedAt` - Last activity time
- Useful for: debugging, SLAs, stale job detection

## Benefits

✅ **Retry Logic** - Jobs automatically retry on failure
✅ **Context Support** - Foundation for cancellation and timeouts
✅ **Error Tracking** - Know why jobs failed
✅ **Audit Trail** - Timestamps for debugging
✅ **Type Safety** - State machine still validates transitions
✅ **Clean API** - Mark methods encapsulate complexity

## What's Next

**Step 10: Thread Safety**
- Add mutex to storage (concurrent access)
- Add mutex to handler registry
- Prepare for worker pool concurrency

Then we'll be ready for async processing with workers!

## Key Takeaway

This step adds **resilience** to our job queue:
- Jobs don't fail permanently on first error
- We track what went wrong and when
- Context enables future cancellation/timeouts
- All while maintaining state machine validation

The retry logic is critical for production systems where temporary failures (network blips, rate limits) are common.
