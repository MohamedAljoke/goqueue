# Step 1: Basic Job Processing

## Goal

Create a simple job that can be processed with a handler function.

## Concepts Learned

- Basic struct creation
- Method receivers in Go
- Simple testing with `testing` package

## Implementation

### goqueue.go

```go
package goqueue

import "fmt"

type Job struct {
    ID      string
    Type    string
    Payload map[string]interface{}
    Status  string
}

func NewJob() *Job {
    job := &Job{
        ID:      "unique-job-id",
        Type:    "default",
        Payload: make(map[string]interface{}),
        Status:  "pending",
    }
    return job
}

func (job *Job) Process() error {
    job.Status = "completed"
    handleJob(job)
    return nil
}

func handleJob(job *Job) error {
    fmt.Printf("processing job %s", job.ID)
    return nil
}
```

## Test

### goqueue_test.go

```go
package goqueue_test

import (
    "testing"
    "github.com/MohamedAljoke/goqueue"
)

func TestGenerateJob(t *testing.T) {
    t.Run("should create job with pending status", func(t *testing.T) {
        job := goqueue.NewJob()
        if job.Status != "pending" {
            t.Errorf("expected job status to be 'pending', got %s", job.Status)
        }
    })
}
```

## Key Points

### Struct Definition
```go
type Job struct {
    ID      string
    Type    string
    Payload map[string]interface{}
    Status  string
}
```
- Exported fields (capitalized) are public
- `map[string]interface{}` allows flexible payload data

### Constructor Pattern
```go
func NewJob() *Job {
    return &Job{...}
}
```
- Returns pointer to avoid copying
- Common Go pattern for initialization

### Method Receiver
```go
func (job *Job) Process() error {
    // job is the receiver
}
```
- `(job *Job)` makes this a method on Job
- Pointer receiver allows modification

## What's Wrong?

### Problems with this implementation:

1. **Hardcoded handler** - Can't have different handlers for different job types
2. **No error handling** - Handler can't fail
3. **No persistence** - Job only exists in memory
4. **Hardcoded ID** - All jobs have same ID

### Next Step

[Step 2: Configurable Handlers](./step2-configurable-handlers.md) - Make handlers flexible
