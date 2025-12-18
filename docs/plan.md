# GoQueue Implementation Plan

## Current Status (Step 10 Complete!)

We've implemented:
- ✅ Clean Architecture (Entity/Storage/Use Case layers)
- ✅ Handler Registry pattern
- ✅ State machine for job status
- ✅ Context support across all layers
- ✅ Job retry fields (Attempts, MaxRetry, CanRetry, BackoffDuration)
- ✅ Mark methods (MarkRunning, MarkCompleted, MarkFailed)
- ✅ Worker pool with goroutines
- ✅ Channels for async job submission
- ✅ Graceful shutdown (context cancellation + WaitGroup)
- ✅ Thread safety for Registry (sync.RWMutex)
- ✅ UpdateJob in storage
- ✅ GetJob API for status checking

## Goal

Build a production-ready asynchronous job queue with worker pool concurrency, matching the architecture of `goqueue_old`.

**Current state:** 90% complete - async worker pool working, but missing critical storage thread safety and actual retry behavior.

---

## Remaining Work (TODO Tomorrow)

### Step 11: Storage Thread Safety (CRITICAL!)

**Problem:** Multiple workers accessing `jobStorage` map concurrently causes data races.

**Test it:**
```bash
go test -race ./...
```

**Expected error:**
```
WARNING: DATA RACE
Write at jobStorage map by goroutine X
Read at jobStorage map by goroutine Y
```

**Solution - Add mutex to MemoryStorage:**

In `internal/storage/job_repository.go`:

```go
type MemoryStorage struct {
    mu   sync.RWMutex                    // NEW: protects jobs map
    jobs map[string]*entity.Job          // NEW: change from global var
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        jobs: make(map[string]*entity.Job),  // Instance-level map
    }
}

func (m *MemoryStorage) SaveJob(ctx context.Context, job *entity.Job) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.jobs[job.ID] = job
}

func (m *MemoryStorage) GetJob(ctx context.Context, id string) (*entity.Job, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    job, exists := m.jobs[id]
    return job, exists
}

func (m *MemoryStorage) UpdateJob(ctx context.Context, job *entity.Job) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if _, exists := m.jobs[job.ID]; !exists {
        return fmt.Errorf("job %s not found", job.ID)
    }
    m.jobs[job.ID] = job
    return nil
}

func (m *MemoryStorage) ClearStorage() {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.jobs = make(map[string]*entity.Job)
}
```

**Key changes:**
1. Add `mu sync.RWMutex` field
2. Change from global `jobStorage` var to instance field `jobs`
3. Lock all operations:
   - `Lock()` for writes (SaveJob, UpdateJob, ClearStorage)
   - `RLock()` for reads (GetJob)

**Verify fix:**
```bash
go test -race ./...  # Should pass!
```

---

### Step 12: Actual Retry Implementation (Enhancement)

**Problem:** We have `CanRetry()` and `BackoffDuration()` methods, but workers don't actually re-queue failed jobs.

**Current behavior:**
- Job fails → MarkFailed() → stays failed (or back to pending, but no retry happens)

**Desired behavior (like goqueue_old):**
- Job fails → MarkFailed() → if CanRetry(), sleep backoff time → re-queue to channel

**Solution - Update worker.processJob():**

In `internal/worker/pool.go`:

```go
func (wp *WorkerPool) worker() {
    defer wp.wg.Done()

    for {
        select {
        case job := <-wp.jobChan:
            handler, err := wp.registry.GetHandler(job.Type)
            if err != nil {
                job.MarkFailed(err)
                wp.storage.SaveJob(wp.ctx, job)
                continue
            }

            // Process job
            err = wp.processor.ProcessJob(wp.ctx, job, handler)

            // If job failed and can retry, re-queue with backoff
            if err != nil && job.CanRetry() {
                go func(j *entity.Job) {
                    backoff := j.BackoffDuration()
                    time.Sleep(backoff)
                    wp.jobChan <- j  // Re-queue
                }(job)
            }

        case <-wp.ctx.Done():
            return
        }
    }
}
```

**Why separate goroutine for retry?**
- Don't block worker during sleep
- Worker can process other jobs immediately
- Backoff happens in background

**Testing retry:**
```go
handler := func(ctx context.Context, payload map[string]interface{}) error {
    return fmt.Errorf("simulated failure")
}

queue.RegisterHandler("retry-test", handler)
queue.Start(3)
job, _ := queue.SubmitJob(ctx, "retry-test", payload, 3)

time.Sleep(20 * time.Second)  // Wait for retries
updated, _ := queue.GetJob(ctx, job.ID)
// updated.Attempts should be 3
// updated.Status should be Failed (after exhausting retries)
```

---

### Step 13: Optional Enhancements (Nice-to-Have)

**Logging (for observability):**
```go
log.Printf("[WORKER-%d] Processing job %s (attempt %d/%d)",
    workerID, job.ID, job.Attempts, job.MaxRetry)
```

**Error constants (better error handling):**
```go
var ErrJobNotFound = errors.New("job not found")
var ErrHandlerNotFound = errors.New("handler not found")
```

**Functional options (flexible config):**
```go
func WithWorkers(count int) Option { ... }
func WithBufferSize(size int) Option { ... }

queue := goqueue.New(
    goqueue.WithWorkers(10),
    goqueue.WithBufferSize(100),
)
```

---

## Summary of Tomorrow's Work

**Priority 1 (MUST DO):**
- ✅ Add storage mutex - fixes critical race condition
- ✅ Test with `go test -race` - verify no races

**Priority 2 (Should Do):**
- Implement actual retry with backoff
- Test retry behavior

**Priority 3 (Nice to Have):**
- Add logging for debugging
- Add error constants
- Functional options pattern

**Documentation:**
- Update step10-worker-pool.md with storage mutex
- Or create step11-storage-thread-safety.md

---

## Completed Steps (Reference)

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
