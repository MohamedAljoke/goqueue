# GoQueue Examples

Real-world examples showing how to use GoQueue in your applications.

## Examples

### 1. Simple Example
**Location:** `examples/simple/`

Basic usage showing:
- Creating a queue
- Registering handlers
- Submitting jobs
- Checking status

```bash
go run examples/simple/main.go
```

**Best for:** Learning the basics, quick prototypes

---

### 2. Web Application Example
**Location:** `examples/webapp/`

HTTP API integration showing:
- REST endpoints for job submission
- Status checking via API
- Graceful shutdown
- Multiple handler types

```bash
go run examples/webapp/main.go

# Test it:
curl -X POST http://localhost:8080/jobs \
  -d '{"type":"send_notification","payload":{"user":"alice"}}'

curl "http://localhost:8080/jobs/status?id=job_123"
```

**Best for:** Web applications, microservices

---

### 3. Service-Oriented Example
**Location:** `examples/service/`

Service pattern integration showing:
- Connecting job handlers to service methods
- Multiple services (Email, Payment, Analytics)
- Namespaced job types
- Real-world architecture

```bash
go run examples/service/main.go
```

**Best for:** Structured applications, service architectures

---

## Usage Pattern Summary

```go
// 1. Create queue
q := goqueue.New(
    goqueue.WithWorkers(5),
    goqueue.WithBufferSize(20),
)

// 2. Register handlers
q.RegisterHandler("job_type", yourHandler)

// 3. Start workers
ctx := context.Background()
go q.Start(ctx)

// 4. Submit jobs
jobID, err := q.Submit(ctx, "job_type", payload, 3)
```

## Next Steps

- Read **[LIBRARY.md](../LIBRARY.md)** for comprehensive API documentation
- Check **[planning/](../planning/)** for production roadmap
- Review handler signature: `func(ctx context.Context, payload map[string]interface{}) error`
