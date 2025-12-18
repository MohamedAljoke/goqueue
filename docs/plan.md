# GoQueue Implementation Plan

## Current Status (Step 8 Complete)

We've implemented:
- ✅ Clean Architecture (Entity/Storage/Use Case layers)
- ✅ Handler Registry pattern
- ✅ Basic job processing (synchronous)
- ✅ State machine for job status
- ✅ Public API (Queue struct)

## Goal

Build a production-ready asynchronous job queue with worker pool concurrency, matching the architecture of `goqueue_old`.

## Pre-Worker Foundation (Steps 9-10)

Before implementing workers, we need thread-safety and retry mechanisms:

### Step 9: Context Support & Job Retry
**Why:** Workers need cancellation, timeouts, and jobs need retry logic

**Add to Job entity (`internal/entity/job.go`):**
- `Attempts int` - Track retry attempts
- `MaxRetry int` - Maximum retry limit
- `Error string` - Store error messages
- `CreatedAt time.Time` - Job creation timestamp
- `UpdatedAt time.Time` - Last update timestamp
- `CanRetry() bool` - Check if job can be retried
- `BackoffDuration() time.Duration` - Exponential backoff calculation

**Update HandlerFunc signature:**
- Change from: `func(payload map[string]interface{}) error`
- Change to: `func(ctx context.Context, payload map[string]interface{}) error`

**Update all layers to accept context:**
- `ProcessJob(ctx context.Context, job *Job, handler HandlerFunc) error`
- `SubmitJob(ctx context.Context, jobType string, payload map[string]interface{}, maxRetry int) (*Job, error)`
- Storage methods: `Save(ctx, job)`, `Get(ctx, id)`, etc.

**Update Job status methods:**
- `MarkRunning()` - Increment attempts, update timestamp
- `MarkCompleted()` - Clear error, update timestamp
- `MarkFailed(err error)` - Store error, decide if pending (retry) or failed (terminal)

### Step 10: Thread Safety & Storage Updates
**Why:** Prevent race conditions when workers access shared resources concurrently

**Add mutex to Handler Registry (`internal/handler/registry.go`):**
```go
type Registry struct {
    mu       sync.RWMutex
    handlers map[string]entity.HandlerFunc
}
```
- Use `RLock/RUnlock` for `GetHandler()`
- Use `Lock/Unlock` for `RegisterHandler()`

**Add mutex to Memory Storage (`internal/storage/job_repository.go`):**
```go
type MemoryStorage struct {
    mu   sync.RWMutex
    jobs map[string]*entity.Job
}
```
- Use `RLock/RUnlock` for reads (`GetJob`)
- Use `Lock/Unlock` for writes (`SaveJob`, `UpdateJob`)

**Extend Storage interface:**
```go
type JobStorage interface {
    SaveJob(ctx context.Context, job *entity.Job) error
    GetJob(ctx context.Context, id string) (*entity.Job, error)
    UpdateJob(ctx context.Context, job *entity.Job) error  // NEW
    ListByStatus(ctx context.Context, status Status) ([]*entity.Job, error)  // Optional
}
```

**Why `UpdateJob` is critical:**
- Workers update job status after processing
- `SaveJob` is for new jobs, `UpdateJob` is for existing jobs
- Can add validation (e.g., job must exist)

## Worker Implementation (Steps 11-12)

### Step 11: Worker Pool
**Why:** Process jobs concurrently with multiple goroutines

**Create `internal/worker/pool.go`:**
- `Pool` struct with job channel, worker count, storage, handlers
- `NewPool(workerCount, bufferSize, storage, handlers)`
- `Enqueue(job)` - Send job to channel
- `Start(ctx)` - Launch worker goroutines
- `worker(ctx, workerID)` - Worker loop: receive jobs, process, handle errors
- `processJob(ctx, workerID, job)` - Execute handler, retry logic, update storage

**Key concepts:**
- Buffered channel: `jobs chan *entity.Job`
- Worker pool pattern with WaitGroup
- Context cancellation for graceful shutdown
- Retry with exponential backoff

### Step 12: Queue Orchestration
**Why:** Connect all pieces into async processing system

**Update `goqueue.go`:**
- Add `pool *worker.Pool` to Queue struct
- `SubmitJob()` → Save to storage + Enqueue to pool
- `Start(ctx)` - Start worker pool (blocking call)
- `GetJob(ctx, id)` - Retrieve job status

**Final public API:**
```go
queue := goqueue.NewQueue()
queue.RegisterHandler("email", emailHandler)
go queue.Start(ctx)  // Run in background
jobID, _ := queue.SubmitJob(ctx, "email", payload, 3)
job, _ := queue.GetJob(ctx, jobID)  // Check status
```

## Optional Enhancements (Step 13+)

- Configuration options (functional options pattern)
- Logging/metrics
- Job cancellation
- Priority queues
- Persistence backends (PostgreSQL, Redis)
- Dead letter queue
- Job scheduling (delayed jobs)

## Testing Strategy

After each step:
1. Update existing tests for new signatures
2. Add new test cases for new functionality
3. Add concurrency tests (race detector: `go test -race`)

## Documentation Strategy

After major milestones:
- Step 9-10 → `step9-context-retry-thread-safety.md`
- Step 11-12 → `step10-workers-async.md`

## Key Differences from goqueue_old

### We improved:
- Clearer separation with `usecase` layer
- More explicit state machine validation
- Better step-by-step learning path

### We kept:
- Same architecture (entity/storage/handler/worker)
- Same patterns (registry, pool, clean architecture)
- Same public API design

## Next Action

Start with **Step 9: Context Support & Job Retry** - foundation for everything else.
