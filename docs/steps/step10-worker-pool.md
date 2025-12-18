# Step 10: Worker Pool & Thread Safety

## Problem

Our current implementation is **synchronous**:
1. `SubmitJob()` blocks until job completes
2. Can only process one job at a time
3. No concurrent processing
4. Slow for multiple jobs

Real-world job queues need:
- **Async submission** - return immediately, process in background
- **Concurrent processing** - multiple jobs running simultaneously
- **Worker pool** - fixed number of workers processing jobs
- **Thread safety** - safe concurrent access to shared data

## Solution

Add a worker pool with goroutines and channels, then discover and fix race conditions.

## Implementation

### 1. Create Worker Pool Structure

**internal/worker/pool.go**

```go
type WorkerPool struct {
    // Dependencies
    storage   storage.JobStorage
    registry  *handler.Registry
    processor *usecase.JobProcessor

    // Communication
    jobChan   chan *entity.Job  // Workers pull jobs from here

    // Lifecycle management
    ctx       context.Context   // For cancellation signal
    cancel    context.CancelFunc
    wg        sync.WaitGroup    // Track active workers
}
```

**Key components:**
- `jobChan` - Channel for job queue (buffer = workerCount × 2)
- `ctx/cancel` - Graceful shutdown signal
- `wg` - Wait for all workers to finish

### 2. Constructor

```go
func NewWorkerPool(storage storage.JobStorage, registry *handler.Registry, workerCount int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    poolChan := make(chan *entity.Job, workerCount*2)

    pool := &WorkerPool{
        storage:   storage,
        registry:  registry,
        processor: usecase.NewJobProcessor(storage),
        jobChan:   poolChan,
        ctx:       ctx,
        cancel:    cancel,
    }

    return pool
}
```

**Buffer size:** `workerCount*2` - good rule of thumb for job queuing.

### 3. Start Workers

```go
func (wp *WorkerPool) Start(workerCount int) {
    for i := 0; i < workerCount; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
}
```

**What it does:**
- Launch `workerCount` goroutines
- Each worker runs concurrently
- `wg.Add(1)` tracks each worker

### 4. Worker Goroutine

```go
func (wp *WorkerPool) worker() {
    defer wp.wg.Done()  // Decrement counter when exiting

    for {
        select {
        case job := <-wp.jobChan:
            // Process job
            handler, err := wp.registry.GetHandler(job.Type)
            if err != nil {
                job.MarkFailed(err)
                wp.storage.SaveJob(wp.ctx, job)
                continue
            }
            wp.processor.ProcessJob(wp.ctx, job, handler)

        case <-wp.ctx.Done():
            // Shutdown signal - exit worker
            return
        }
    }
}
```

**How it works:**
- Infinite loop waiting on channels
- `select` waits for either:
  - New job from `jobChan`
  - Cancellation signal from `ctx.Done()`
- Process job → get handler → call processor
- Exit on shutdown signal

### 5. Submit and Stop Methods

```go
func (wp *WorkerPool) Submit(job *entity.Job) {
    wp.jobChan <- job
}

func (wp *WorkerPool) Stop() {
    wp.cancel()        // Signal all workers to stop
    close(wp.jobChan)  // No more jobs accepted
    wp.wg.Wait()       // Wait for all workers to exit
}
```

**Shutdown flow:**
1. `cancel()` closes `ctx.Done()` channel
2. `close(jobChan)` prevents new submissions
3. `wg.Wait()` blocks until all workers finish current jobs

### 6. Update Storage - Add UpdateJob

**internal/storage/job_repository.go**

```go
type JobStorage interface {
    SaveJob(ctx context.Context, job *entity.Job)
    GetJob(ctx context.Context, id string) (*entity.Job, bool)
    UpdateJob(ctx context.Context, job *entity.Job) error  // NEW
    ClearStorage()
}

func (m *MemoryStorage) UpdateJob(ctx context.Context, job *entity.Job) error {
    if _, exists := jobStorage[job.ID]; !exists {
        return fmt.Errorf("job %s not found", job.ID)
    }
    jobStorage[job.ID] = job
    return nil
}
```

**Why separate methods?**
- `SaveJob()` - Create new job
- `UpdateJob()` - Update existing (returns error if not found)

### 7. Update Queue - Make it Async

**Remove processor, add pool:**
```go
type Queue struct {
    storage  storage.JobStorage
    registry *handler.Registry
    pool     *worker.WorkerPool  // NEW: worker pool
}
```

**Update constructor:**
```go
func NewQueue() *Queue {
    storage := storage.NewMemoryStorage()
    registry := handler.NewRegistry()
    workerCount := 5
    pool := worker.NewWorkerPool(storage, registry, workerCount)

    return &Queue{
        storage:  storage,
        registry: registry,
        pool:     pool,
    }
}
```

**Add lifecycle methods:**
```go
func (q *Queue) Start(workerCount int) {
    q.pool.Start(workerCount)
}

func (q *Queue) Stop() {
    q.pool.Stop()
}
```

**Update SubmitJob - Now Async:**
```go
func (q *Queue) SubmitJob(ctx context.Context, jobType string, payload map[string]interface{}, maxRetry int) (*entity.Job, error) {
    // Validate handler exists
    _, err := q.registry.GetHandler(jobType)
    if err != nil {
        return nil, err
    }

    // Create job
    job := entity.NewJob(maxRetry)
    job.Type = jobType
    job.Payload = payload

    // Save to storage
    q.storage.SaveJob(ctx, job)

    // Submit to worker pool (async)
    q.pool.Submit(job)

    // Return immediately - job will be processed in background
    return job, nil
}
```

**Add GetJob - Check Status:**
```go
func (q *Queue) GetJob(ctx context.Context, jobID string) (*entity.Job, error) {
    job, exists := q.storage.GetJob(ctx, jobID)
    if !exists {
        return nil, fmt.Errorf("job %s not found", jobID)
    }
    return job, nil
}
```

**Usage flow:**
```go
queue := goqueue.NewQueue()
queue.RegisterHandler("email", emailHandler)
queue.Start(5)  // Start 5 workers

// Submit returns immediately
job, _ := queue.SubmitJob(ctx, "email", payload, 3)
// job.Status = Pending

// Later check status
updated, _ := queue.GetJob(ctx, job.ID)
// updated.Status might be Processing, Completed, or Failed

queue.Stop()  // Graceful shutdown
```

## Discovering Race Conditions

### The Problem

Run existing tests with race detector:
```bash
go test -race ./...
```

**You'll see:**
```
==================
WARNING: DATA RACE
Read at 0x... by goroutine 10:
  github.com/MohamedAljoke/goqueue/internal/handler.(*Registry).GetHandler()

Previous write at 0x... by goroutine 12:
  github.com/MohamedAljoke/goqueue/internal/handler.(*Registry).RegisterHandler()
==================
```

**Why?** Go maps are **not thread-safe**:
- Multiple workers calling `GetHandler()` (reading)
- Potentially calling `RegisterHandler()` while workers run (writing)
- Concurrent map access = undefined behavior

### The Solution - Add Mutex

**internal/handler/registry.go**

```go
type Registry struct {
    handlers map[string]entity.HandlerFunc
    mu       sync.RWMutex  // NEW: protects handlers map
}
```

**Lock writes:**
```go
func (r *Registry) RegisterHandler(jobType string, handler entity.HandlerFunc) {
    r.mu.Lock()         // Write lock - exclusive
    defer r.mu.Unlock()
    r.handlers[jobType] = handler
}
```

**Lock reads:**
```go
func (r *Registry) GetHandler(jobType string) (entity.HandlerFunc, error) {
    r.mu.RLock()        // Read lock - allows concurrent reads
    defer r.mu.RUnlock()

    handler, exists := r.handlers[jobType]
    if !exists {
        return nil, fmt.Errorf("no handler registered for job type: %s", jobType)
    }
    return handler, nil
}
```

**Run tests again:**
```bash
go test -race ./...
```

✅ **No race conditions detected!**

### Demonstrating the Race (For Students)

**Create internal/handler/registry_race_test.go:**

```go
func TestRegistryConcurrentAccess(t *testing.T) {
    registry := handler.NewRegistry()

    registry.RegisterHandler("email", func(ctx context.Context, payload map[string]interface{}) error {
        return nil
    })

    var wg sync.WaitGroup

    // 10 workers reading concurrently
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                registry.GetHandler("email")
            }
        }()
    }

    // 5 goroutines writing concurrently
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 50; j++ {
                registry.RegisterHandler("sms", func(ctx context.Context, payload map[string]interface{}) error {
                    return nil
                })
            }
        }(i)
    }

    wg.Wait()
}
```

**Run with `-race` → shows race condition before mutex is added!**

## What We Learned

### Go Concurrency Concepts

**1. Channels**
```go
jobChan := make(chan *entity.Job, bufferSize)
jobChan <- job              // Send to channel
job := <-jobChan            // Receive from channel
close(jobChan)              // Close channel
```

**2. Goroutines**
```go
go wp.worker()  // Launch concurrent function
```

**3. Select Statement**
```go
select {
case job := <-jobChan:
    // Process job
case <-ctx.Done():
    // Shutdown
}
```

**4. Context Cancellation**
```go
ctx, cancel := context.WithCancel(context.Background())
cancel()            // Signal cancellation
<-ctx.Done()        // Receive signal
```

**5. WaitGroup**
```go
var wg sync.WaitGroup
wg.Add(1)      // Increment counter
go func() {
    defer wg.Done()  // Decrement when done
}()
wg.Wait()      // Block until counter = 0
```

**6. RWMutex (Read-Write Mutex)**
```go
var mu sync.RWMutex
mu.Lock()      // Exclusive write lock
mu.Unlock()
mu.RLock()     // Shared read lock (multiple readers OK)
mu.RUnlock()
```

### Worker Pool Pattern

```
NewQueue()
    ↓
NewWorkerPool() → create channel, context
    ↓
Start(5) → launch 5 worker goroutines
    ↓
Workers loop: select { jobChan | ctx.Done() }
    ↓
SubmitJob() → create job → save → push to channel → return
    ↓
Worker picks up → get handler → process → save
    ↓
GetJob() → check status
    ↓
Stop() → cancel() → close(channel) → wg.Wait()
```

### Why RWMutex vs Mutex?

**Regular Mutex:**
- Only one goroutine can hold lock (read OR write)
- Simple but slower for read-heavy workloads

**RWMutex:**
- Multiple readers can hold `RLock()` simultaneously
- Only one writer with `Lock()` (exclusive)
- Better for registry (many reads, few writes)

### Race Detector

```bash
go test -race          # Run tests with race detection
go build -race         # Build with race detection
go run -race main.go   # Run with race detection
```

**Detects:**
- Concurrent map access
- Unsynchronized variable access
- Data races across goroutines

**Production:** Don't ship with `-race` (adds overhead), but use it in testing!

## Benefits

✅ **Async Processing** - SubmitJob returns immediately
✅ **Concurrent Execution** - Multiple jobs processed simultaneously
✅ **Worker Pool Pattern** - Fixed number of workers, prevents resource exhaustion
✅ **Graceful Shutdown** - Finish current jobs before stopping
✅ **Thread Safety** - Mutex protects shared data
✅ **Channels** - Safe communication between goroutines
✅ **Context Cancellation** - Clean shutdown signal propagation

## What's Next

**Step 11: Retry with Backoff** (optional enhancement)
- Implement actual retry logic with `BackoffDuration()`
- Use `time.Sleep()` in worker before retrying
- Show failed → pending → processing flow

**Step 12: Persistent Storage** (optional)
- Replace in-memory storage with database
- Add mutex to storage as well
- Transaction safety

**Or we're done!** The queue is now production-ready for basic use cases.

## Key Takeaways

This step adds **concurrency** to our job queue:

1. **Worker Pool** - Multiple workers processing jobs concurrently
2. **Channels** - Safe communication between Queue and Workers
3. **Thread Safety** - Discovered race conditions organically, fixed with mutex
4. **Async API** - Submit returns immediately, check status later
5. **Graceful Shutdown** - Clean worker termination

**Teaching approach:**
- Built worker pool first
- Ran race detector
- Discovered the problem (data race)
- Fixed with mutex

This shows students **why** thread safety matters, not just **what** to do!
