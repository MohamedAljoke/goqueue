# Step 6: State Machine Pattern

## Problem

- Status was just a string - no type safety
- No validation for invalid state transitions
- Job was tightly coupled with storage

## Solution

- Create Status type with constants
- Add state transition validation
- Prepare to decouple storage (next step)

## Concepts Learned

- Type aliases for type safety
- State machine pattern
- Validation logic
- Terminal states concept

## Implementation

### Define Status Type

```go
type Status string

const (
    StatusPending    Status = "pending"
    StatusProcessing Status = "processing"
    StatusCompleted  Status = "completed"
    StatusFailed     Status = "failed"
)
```

### Update Job Struct

```go
type Job struct {
    ID      string
    Type    string
    Payload map[string]interface{}
    Status  Status              // Changed from string
    Error   string              // NEW: Store error messages
}
```

### State Transition Validation

```go
func (job *Job) canTransition(from, to Status) bool {
    validTransitions := map[Status][]Status{
        StatusPending: {
            StatusProcessing,
        },
        StatusProcessing: {
            StatusCompleted,
            StatusFailed,
        },
        StatusFailed:    {},  // Terminal state - no exits
        StatusCompleted: {},  // Terminal state - no exits
    }

    for _, allowed := range validTransitions[from] {
        if allowed == to {
            return true
        }
    }
    return false
}
```

### Status Change Method

```go
func (job *Job) ChangeStatus(status Status) error {
    if !job.canTransition(job.Status, status) {
        return fmt.Errorf(
            "invalid job status transition: %s -> %s",
            job.Status,
            status,
        )
    }

    job.Status = status
    SaveJob(job)  // Still coupled - will fix next

    return nil
}
```

### Updated Process Method

```go
func (job *Job) Process(handler HandlerFunc) error {
    // Use ChangeStatus for validation
    if err := job.ChangeStatus(StatusProcessing); err != nil {
        return err
    }

    err := handler(job.Payload)
    if err != nil {
        job.Error = err.Error()  // Store error message
        _ = job.ChangeStatus(StatusFailed)
        return fmt.Errorf("error handling process: %w", err)
    }

    if err := job.ChangeStatus(StatusCompleted); err != nil {
        return err
    }

    return nil
}
```

## Key Concepts

### Type Aliases

```go
type Status string

var s Status = "pending"       // OK
var s Status = StatusPending   // OK (type-safe)
var s string = StatusPending   // OK (underlying type)
var s Status = 123            // ERROR (type safety!)
```

Type alias provides:
- Type safety
- Better documentation
- IDE autocomplete

### State Machine

```
┌─────────┐
│ Pending │
└────┬────┘
     │
     ▼
┌────────────┐
│ Processing │
└─────┬──────┘
      │
      ├──────────┐
      │          │
      ▼          ▼
┌───────────┐  ┌────────┐
│ Completed │  │ Failed │
└───────────┘  └────────┘
```

Valid transitions:
- pending → processing ✅
- processing → completed ✅
- processing → failed ✅
- completed → processing ❌
- failed → pending ❌

### Terminal States

```go
StatusFailed:    {},  // Empty slice = no valid transitions
StatusCompleted: {},  // Empty slice = no valid transitions
```

Once a job reaches these states, it cannot transition further (unless you implement retry logic).

### Validation Pattern

```go
func (job *Job) ChangeStatus(status Status) error {
    if !job.canTransition(job.Status, status) {
        return fmt.Errorf("invalid transition")
    }

    job.Status = status
    return nil
}
```

Prevents invalid state changes:
```go
job.Status = StatusCompleted
job.ChangeStatus(StatusProcessing)  // ERROR: Invalid transition
```

## State Transition Table

| From | To | Valid? |
|------|-----|--------|
| Pending | Processing | ✅ |
| Pending | Completed | ❌ |
| Pending | Failed | ❌ |
| Processing | Completed | ✅ |
| Processing | Failed | ✅ |
| Processing | Pending | ❌ |
| Completed | * | ❌ |
| Failed | * | ❌ |

## Error Storage

```go
type Job struct {
    // ...
    Error string  // Store error message when job fails
}

// In Process:
if err != nil {
    job.Error = err.Error()  // Save error message
    job.ChangeStatus(StatusFailed)
}
```

Now you can inspect why a job failed:
```go
job, _ := GetJob("job_123")
if job.Status == StatusFailed {
    fmt.Println("Failed because:", job.Error)
}
```

## Type Safety Benefits

### Before (String)
```go
job.Status = "procesing"      // Typo! No error
job.Status = "PENDING"        // Wrong case! No error
job.Status = "completed"      // OK
```

### After (Status Type)
```go
job.Status = StatusProcessing // OK - autocomplete helps
job.Status = StatusPENDING    // Compile error if doesn't exist
job.Status = "completed"      // Still works (underlying type)
job.ChangeStatus(StatusPending) // Runtime error - validates transition
```

## Testing State Transitions

```go
func TestInvalidTransition(t *testing.T) {
    job := internal.NewJob()
    job.Status = internal.StatusCompleted

    err := job.ChangeStatus(internal.StatusProcessing)
    if err == nil {
        t.Errorf("expected error for invalid transition")
    }
}
```

## What's Still Wrong?

1. **Still coupled to storage** - ChangeStatus() calls SaveJob()
2. **Generic ChangeStatus** - Could be more specific (MarkCompleted, MarkFailed, etc.)
3. **No retry logic** - Failed jobs can't retry
4. **No attempt counter** - Can't limit retries

## Next Steps

### Step 7: Decouple Storage
- Remove SaveJob() calls from job.go
- Make caller responsible for persistence
- Pure domain logic

### Step 8: Specific Transition Methods
```go
func (job *Job) MarkRunning()
func (job *Job) MarkCompleted()
func (job *Job) MarkFailed(err error)
```

### Step 9: Add Retry Logic
- Attempt counter
- Max retry limit
- Exponential backoff

## Summary

State machine pattern provides:
✅ Type safety
✅ Invalid transition prevention
✅ Clear state flow
✅ Self-documenting code
✅ Error storage
