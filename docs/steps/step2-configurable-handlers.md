# Step 2: Configurable Handlers

## Problem

The handler was hardcoded - we needed different handlers for different job types.

## Solution

Introduce `HandlerFunc` type and pass it to `Process()`.

## Concepts Learned

- First-class functions in Go
- Function types
- Dependency injection pattern

## Implementation

### Define HandlerFunc Type

```go
type HandlerFunc func(payload map[string]interface{}) error
```

This creates a **function type** - functions matching this signature can be assigned to variables of type `HandlerFunc`.

### Update Process Method

```go
func (job *Job) Process(handler HandlerFunc) error {
    job.Status = "completed"
    handler(job.Payload)
    return nil
}
```

Now `Process()` accepts any function that matches the `HandlerFunc` signature.

## Test

```go
t.Run("should have status completed after processing", func(t *testing.T) {
    job := goqueue.NewJob()

    // Define a custom handler
    handler := func(payload map[string]interface{}) error {
        return nil
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

## Key Concepts

### First-Class Functions

In Go, functions are first-class citizens:
```go
// Function type
type HandlerFunc func(payload map[string]interface{}) error

// Function variable
var handler HandlerFunc = func(payload map[string]interface{}) error {
    return nil
}

// Pass function as parameter
job.Process(handler)
```

### Anonymous Functions

```go
handler := func(payload map[string]interface{}) error {
    fmt.Println("Processing...")
    return nil
}
```
- Function without a name
- Can capture variables from surrounding scope (closure)

### Dependency Injection

Instead of hardcoding `handleJob()`, we **inject** the handler:
```go
// Before: Hardcoded
func (job *Job) Process() {
    handleJob(job)  // Fixed dependency
}

// After: Injected
func (job *Job) Process(handler HandlerFunc) {
    handler(job.Payload)  // Flexible dependency
}
```

## Usage Example

```go
// Email handler
emailHandler := func(payload map[string]interface{}) error {
    email := payload["email"].(string)
    fmt.Printf("Sending email to %s\n", email)
    return nil
}

// Payment handler
paymentHandler := func(payload map[string]interface{}) error {
    amount := payload["amount"].(float64)
    fmt.Printf("Processing payment of $%.2f\n", amount)
    return nil
}

// Use different handlers
job1.Process(emailHandler)
job2.Process(paymentHandler)
```

## What's Still Wrong?

1. **No error handling** - We call `handler()` but ignore its return value
2. **Status always "completed"** - Even if handler fails
3. **No persistence** - Jobs lost when program ends

## Next Step

[Step 3: Error Handling](./step3-error-handling.md) - Handle failures properly
