# Technical Dependencies & Prerequisites - Go Job Processor

## Overview
This document outlines the technical dependencies, prerequisites, and architectural decisions needed before development begins.

---

## Development Environment

### Required Tools
- **Go:** Version 1.21 or higher
- **Git:** For version control
- **Docker:** Version 20.10+ (for local PostgreSQL and deployment)
- **Docker Compose:** For local development environment
- **Make:** For build automation (optional but recommended)
- **golangci-lint:** For code linting

### Recommended Tools
- **Air:** For hot reload during development
- **GoLand/VSCode:** IDE with Go support
- **Postman/Insomnia:** For API testing
- **pgAdmin/TablePlus:** For database inspection

---

## Go Dependencies

### Core Dependencies

#### Standard Library
- `context` - Context management and cancellation
- `sync` - Concurrency primitives (Mutex, WaitGroup)
- `time` - Time operations and timers
- `net/http` - HTTP server
- `encoding/json` - JSON serialization

#### External Packages (Initial)
```go
// UUID generation
github.com/google/uuid v1.5.0

// Database driver (when ready)
github.com/jackc/pgx/v5 v5.5.0

// Database migrations
github.com/golang-migrate/migrate/v4 v4.17.0

// Configuration management
github.com/spf13/viper v1.18.0
// OR
github.com/kelseyhightower/envconfig v1.4.0

// Structured logging
log/slog (standard library, Go 1.21+)
// OR
github.com/rs/zerolog v1.31.0

// HTTP router (optional, can use net/http)
github.com/gorilla/mux v1.8.1
// OR
github.com/gin-gonic/gin v1.9.1

// Testing
github.com/stretchr/testify v1.8.4
github.com/testcontainers/testcontainers-go v0.27.0
```

### Metrics & Observability (Optional)
```go
// Prometheus metrics
github.com/prometheus/client_golang v1.18.0

// Distributed tracing (advanced)
go.opentelemetry.io/otel v1.21.0
```

---

## Project Structure

### Recommended Clean Architecture Layout

```
go-job-processor/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/                     # Private application code
│   ├── domain/                  # Domain entities and business logic
│   │   ├── job.go              # Job entity
│   │   ├── job_test.go
│   │   └── repository.go       # Repository interface
│   ├── application/             # Application services and use cases
│   │   ├── service/
│   │   │   ├── job_service.go  # Job business logic
│   │   │   └── job_service_test.go
│   │   └── usecase/
│   │       ├── create_job.go
│   │       └── cancel_job.go
│   ├── infrastructure/          # External implementations
│   │   ├── memory/
│   │   │   └── repository.go   # In-memory repository
│   │   ├── postgres/
│   │   │   └── repository.go   # PostgreSQL repository
│   │   └── config/
│   │       └── config.go       # Configuration loading
│   ├── worker/                  # Worker pool and execution
│   │   ├── pool.go
│   │   ├── worker.go
│   │   └── pool_test.go
│   ├── dispatcher/              # Job dispatcher
│   │   ├── dispatcher.go
│   │   └── dispatcher_test.go
│   ├── handler/                 # Job handlers
│   │   ├── handler.go          # Handler interface
│   │   ├── registry.go         # Handler registry
│   │   ├── email_handler.go
│   │   ├── report_handler.go
│   │   └── file_handler.go
│   ├── retry/                   # Retry logic
│   │   ├── backoff.go
│   │   └── scheduler.go
│   └── api/                     # HTTP API layer
│       ├── server.go
│       ├── handlers/
│       │   ├── job_handler.go
│       │   └── health_handler.go
│       ├── middleware/
│       │   ├── logging.go
│       │   └── error_handler.go
│       └── dto/
│           ├── request.go
│           └── response.go
├── pkg/                         # Public libraries (if any)
│   └── logger/
│       └── logger.go
├── migrations/                  # Database migrations
│   ├── 000001_create_jobs_table.up.sql
│   └── 000001_create_jobs_table.down.sql
├── scripts/                     # Build and utility scripts
│   ├── run-tests.sh
│   └── run-migrations.sh
├── docs/                        # Documentation
│   ├── architecture.md
│   └── api/
│       └── openapi.yaml
├── docker/
│   └── Dockerfile
├── .env.example                 # Example environment variables
├── .gitignore
├── docker-compose.yml
├── Makefile                     # Build automation
├── go.mod
├── go.sum
└── README.md
```

---

## Architecture Decisions

### ADR-001: Use Clean Architecture
**Status:** Accepted

**Context:**
Need clear separation of concerns for testability and maintainability.

**Decision:**
Follow Clean Architecture with domain, application, and infrastructure layers.

**Consequences:**
- More boilerplate initially
- Easier to test
- Easier to swap implementations (in-memory → database)

---

### ADR-002: Use Channels for Job Distribution
**Status:** Accepted

**Context:**
Need a way to distribute jobs to workers safely.

**Decision:**
Use Go channels with a worker pool pattern.

**Consequences:**
- Provides natural backpressure
- Simple coordination
- Limited to single process (distributed queues need different approach)

---

### ADR-003: Use Context for Cancellation
**Status:** Accepted

**Context:**
Need a way to cancel running jobs gracefully.

**Decision:**
Pass context.Context to all job handlers and use context cancellation.

**Consequences:**
- Idiomatic Go
- Handlers must respect context
- Enables timeout support

---

### ADR-004: Start with In-Memory, Then Database
**Status:** Accepted

**Context:**
Need to validate core logic before adding database complexity.

**Decision:**
Start with in-memory repository, add PostgreSQL later.

**Consequences:**
- Faster initial development
- Repository interface must be well-designed
- Migration path needed

---

### ADR-005: Use Repository Pattern
**Status:** Accepted

**Context:**
Need abstraction for data persistence.

**Decision:**
Define Repository interface in domain layer, implement in infrastructure.

**Consequences:**
- Easy to swap implementations
- Testable with mock repositories
- Clear dependency direction (infrastructure depends on domain)

---

### ADR-006: Exponential Backoff for Retries
**Status:** Accepted

**Context:**
Need retry strategy that doesn't overwhelm downstream services.

**Decision:**
Use exponential backoff with optional jitter.

**Consequences:**
- Reduces load on failing services
- May delay job completion
- Need to make configurable

---

### ADR-007: Handler Interface for Extensibility
**Status:** Accepted

**Context:**
Different job types need different execution logic.

**Decision:**
Define Handler interface with Execute method, use registry for lookup.

**Consequences:**
- Easy to add new job types
- Handlers are independently testable
- Requires registration boilerplate

---

### ADR-008: PostgreSQL for Persistence
**Status:** Accepted

**Context:**
Need durable storage for production use.

**Decision:**
Use PostgreSQL with pgx driver.

**Consequences:**
- Reliable and well-supported
- Requires database setup
- Use SKIP LOCKED for job claiming

---

### ADR-009: Standard Library HTTP Server
**Status:** Proposed

**Context:**
Need HTTP API but want to minimize dependencies initially.

**Decision:**
Start with net/http, consider gin/echo if routing becomes complex.

**Consequences:**
- No external dependencies initially
- More manual routing
- Can switch later if needed

---

### ADR-010: Structured Logging with slog
**Status:** Accepted

**Context:**
Need logging for production debugging and monitoring.

**Decision:**
Use Go 1.21+ slog for structured logging in JSON format.

**Consequences:**
- Standard library (no dependency)
- Structured logs easier to parse
- Requires Go 1.21+

---

## Story Dependencies Graph

### Critical Path (Must Complete in Order)

```
Milestone 1: Foundation
  STORY-001 (Job Entity)
    ↓
  STORY-002 (Job Status)
    ↓
  STORY-003 (State Machine)
    ↓
  STORY-004 (Business Logic)
    ↓
  STORY-005 (Repository Interface)
    ↓
  STORY-006 (In-Memory Repository)
    ↓
  STORY-007 (Repository Queries)

Milestone 2: Concurrency
  STORY-008 (Worker Pool)
    ↓
  STORY-009 (Worker Loop)
    ↓
  STORY-010 (Dispatcher)
    ↓
  STORY-013 (Status Update Service)
    ↓
  STORY-014 (Job Execution Workflow)
    ↓
  STORY-016 (Integration)
    ↓
  STORY-012 (Concurrent Tests)

Milestone 3: Retry
  STORY-017 (Backoff Strategy)
    ↓
  STORY-018 (Retry Logic)
    ↓
  STORY-019 (Retry Scheduler)
    ↓
  STORY-020 (Retry Tests)
```

### Parallel Work Opportunities

**Can work in parallel:**
- STORY-011 (Worker Config) alongside STORY-010 (Dispatcher)
- STORY-015 (Panic Recovery) alongside STORY-014 (Job Execution)
- Epic 6 (Cancellation) can start after STORY-016
- Epic 7 (Handlers) can start after STORY-014
- Epic 10 (Database) can start after STORY-005 (Repository Interface)

---

## Testing Strategy

### Unit Tests
- All domain logic (state machine, business rules)
- Individual components (handlers, backoff strategy)
- Target: >80% coverage for domain and application layers

### Integration Tests
- Worker pool with real concurrency
- API endpoints with test repository
- Database repository with test containers

### Race Condition Testing
- Run all tests with `-race` flag
- Verify no data races in concurrent execution
- Test with high concurrency (100+ jobs)

### Load Testing (Optional)
- Benchmark with 1000+ jobs
- Measure throughput (jobs/second)
- Identify bottlenecks

---

## Database Setup

### Local Development
```bash
# Using Docker Compose
docker-compose up -d postgres

# Connection string
postgresql://jobprocessor:password@localhost:5432/jobprocessor?sslmode=disable
```

### Schema Requirements
```sql
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);
CREATE INDEX idx_jobs_status_created_at ON jobs(status, created_at);
```

---

## Configuration Management

### Environment Variables
```bash
# Server
SERVER_PORT=8080

# Worker Pool
WORKER_POOL_SIZE=10
JOB_QUEUE_SIZE=100

# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/dbname

# Retry
MAX_RETRY_ATTEMPTS=3
RETRY_BASE_DELAY=1s
RETRY_MAX_DELAY=60s

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Shutdown
SHUTDOWN_TIMEOUT=30s
```

---

## Development Workflow

### Phase 1: Foundation (Weeks 1-2)
1. Setup project structure
2. Implement domain model
3. Create in-memory repository
4. Write unit tests

### Phase 2: Concurrency (Weeks 3-4)
1. Implement worker pool
2. Create dispatcher
3. Integrate job execution
4. Test with `-race` flag

### Phase 3: Features (Weeks 5-7)
1. Add retry logic
2. Implement cancellation
3. Create handler system

### Phase 4: API (Weeks 8-9)
1. Build HTTP server
2. Implement endpoints
3. Add validation
4. Write API tests

### Phase 5: Production (Weeks 10-14)
1. Add observability
2. Implement database persistence
3. Create Docker images
4. Write documentation

---

## Pre-Development Checklist

### Before Starting Development
- [ ] Go 1.21+ installed
- [ ] Docker and Docker Compose installed
- [ ] Project repository created
- [ ] Project structure scaffolded
- [ ] `.gitignore` configured
- [ ] `go.mod` initialized
- [ ] Development environment documented
- [ ] Team aligned on architecture decisions

### Before Starting Each Epic
- [ ] All prerequisite stories completed
- [ ] Dependencies available
- [ ] Acceptance criteria clear
- [ ] Test strategy defined

### Before Milestone Completion
- [ ] All stories in milestone completed
- [ ] Unit tests passing
- [ ] Integration tests passing
- [ ] No race conditions detected
- [ ] Code reviewed
- [ ] Documentation updated

---

## Risk Mitigation

### Technical Risks

#### Risk: Race Conditions in Concurrent Code
**Mitigation:**
- Run all tests with `-race` flag
- Use proper synchronization primitives
- Code review focusing on concurrency

#### Risk: Database Migration Issues
**Mitigation:**
- Use proven migration tool (golang-migrate)
- Test migrations in dev environment
- Keep migrations reversible
- Version control migrations

#### Risk: Context Cancellation Edge Cases
**Mitigation:**
- Extensive testing of cancellation scenarios
- Document context handling patterns
- Review handler implementations

#### Risk: Performance Bottlenecks
**Mitigation:**
- Benchmark early and often
- Profile with pprof
- Load test with realistic data
- Monitor metrics in production

---

## External Service Dependencies

### Development
- **PostgreSQL:** Local via Docker Compose
- **None required initially:** System can run standalone

### Production (Future)
- **PostgreSQL:** Managed database (AWS RDS, Google Cloud SQL, etc.)
- **Monitoring:** Prometheus + Grafana (optional)
- **Tracing:** Jaeger/Zipkin (optional)
- **Logging:** Centralized logging (ELK, Datadog, etc.)

---

## Success Criteria

### Milestone 1 Success
- Domain model complete and tested
- In-memory repository working
- >80% test coverage on domain layer

### Milestone 2 Success
- Jobs execute concurrently without errors
- No race conditions detected
- Workers start and stop cleanly

### MVP Success (Milestone 6)
- All core capabilities working:
  - Create jobs via API
  - Jobs execute asynchronously
  - Status tracked correctly
  - Failed jobs retry automatically
  - Jobs can be canceled
  - Status queryable via API
- API documented
- Integration tests passing

### Production-Ready Success (Milestone 9)
- Database persistence working
- Observability in place
- Graceful shutdown implemented
- Docker images built
- Deployment documented

---

## Questions to Resolve Before Development

1. **Authentication/Authorization:** Is API authentication needed? (Recommendation: Start without, add later if needed)
2. **Rate Limiting:** Should the API have rate limits? (Recommendation: Not for MVP)
3. **Multi-tenancy:** Support multiple clients/tenants? (Recommendation: Not for MVP)
4. **Distributed Deployment:** Multi-instance support from start? (Recommendation: Design for it, implement in M8)
5. **Job Priority:** Do some jobs have higher priority? (Recommendation: Not for MVP, can add later)
6. **Job Scheduling:** Schedule jobs for future execution? (Recommendation: Not for MVP)
7. **Dead Letter Queue:** What happens to jobs that fail max attempts? (Recommendation: Keep in DB with failed status)
8. **Monitoring:** Which metrics are most important? (Recommendation: Start with basic counts, expand later)

---

## Next Steps

1. **Review this plan** with the team
2. **Resolve open questions** (see above)
3. **Setup development environment** (Docker, Go, etc.)
4. **Create project repository** and initialize structure
5. **Begin STORY-001** (Define Job Entity)
6. **Hold daily standups** to track progress
7. **Review after each milestone** to adjust plan if needed
