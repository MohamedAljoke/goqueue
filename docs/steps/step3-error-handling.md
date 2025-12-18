# Step 3: Error Handling

## Problem

Handlers could fail, but we weren't handling errors properly.

## Solution

- Check handler return value
- Set status to "failed" on error
- Add "processing" status to show job is running

## Concepts Learned

- Error handling patterns in Go
- Error wrapping with `%w`
- Testing error cases

## Implementation

### Updated Process Method

```go
func (job *Job) Process(handler HandlerFunc) error {
    job.Status = "processing"  // Mark as running

    err := handler(job.Payload)
    if err != nil {
        job.Status = "failed"
        return fmt.Errorf("error handling process %w", err)
    }

    job.Status = "completed"
    return nil
}
```

### Status Flow

```
pending → processing → completed
                    → failed
```

## Tests

### Success Case

```go
t.Run("should have status completed after processing", func(t *testing.T) {
    job := goqueue.NewJob()
    handler := func(payload map[string]interface{}) error {
        return nil  // Success
    }

    err := job.Process(handler)
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if job.Status != "completed" {
        t.Errorf("expected job status to be 'completed', got %s", job.Status)
    }
})
```

### Failure Case

```go
t.Run("should have status failed after processing with error", func(t *testing.T) {
    job := goqueue.NewJob()
    handler := func(payload map[string]interface{}) error {
        return fmt.Errorf("handler error")  // Failure
    }

    err := job.Process(handler)
    if err == nil {
        t.Errorf("expected error, got nil")
    }
    if job.Status != "failed" {
        t.Errorf("expected job status to be 'failed', got %s", job.Status)
    }
})
```

## Key Concepts

### Error Handling in Go

```go
result, err := someFunction()
if err != nil {
    // Handle error
    return err
}
// Use result
```

Go uses explicit error returns instead of exceptions.

### Error Wrapping

```go
return fmt.Errorf("error handling process %w", err)
```

- `%w` wraps the original error
- Preserves error chain
- Can be unwrapped later with `errors.Unwrap()`

**Example:**
```go
originalErr := fmt.Errorf("database connection failed")
wrappedErr := fmt.Errorf("failed to save user: %w", originalErr)

// Error message: "failed to save user: database connection failed"
```

### Error Creation

```go
// Simple error
err := fmt.Errorf("something went wrong")

// Formatted error
err := fmt.Errorf("invalid status: %s", status)

// Wrapped error
err := fmt.Errorf("operation failed: %w", originalErr)
```

## Testing Errors

Always test both success and failure paths:

```go
// Test success
err := job.Process(successHandler)
if err != nil {
    t.Errorf("expected no error")
}

// Test failure
err = job.Process(failureHandler)
if err == nil {
    t.Errorf("expected error")
}
```

## Status Tracking

Adding "processing" status gives us visibility:

```
pending     - Job created, waiting
processing  - Handler is running
completed   - Handler succeeded
failed      - Handler returned error
```

## What's Still Missing?

1. **No persistence** - Status changes are lost
2. **No error message storage** - Can't see why job failed
3. **No retry logic** - Failed jobs can't retry
4. **String-based status** - No type safety

## Next Step

[Step 4: Storage Layer](./step4-add-storage.md) - Persist job state
