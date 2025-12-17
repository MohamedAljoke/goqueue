# User Stories & Tasks - Go Job Processor

## Overview
This document contains detailed user stories and tasks for each epic. These can be directly converted to Jira tickets.

**Story Point Scale:**
- 1 point: ~1-2 hours
- 2 points: ~3-4 hours
- 3 points: ~1 day
- 5 points: ~2-3 days
- 8 points: ~1 week
- 13 points: ~2 weeks

---

## Epic 1: Core Domain Model

### STORY-001: Define Job Entity
**Story Points:** 3
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need a Job entity that represents all aspects of a job so that the system can track job lifecycle.

**Acceptance Criteria:**
- [ ] Job struct includes: ID, Type, Payload, Status, Attempts, MaxAttempts, LastError, CreatedAt, UpdatedAt
- [ ] Job ID is a UUID
- [ ] Job Type is a string type alias
- [ ] Payload is stored as []byte
- [ ] Timestamps use time.Time
- [ ] All fields properly documented

**Technical Tasks:**
- [ ] Create `domain/job.go`
- [ ] Define Job struct with all attributes
- [ ] Add constructor function NewJob()
- [ ] Add godoc comments
- [ ] Write unit tests for Job creation

**Dependencies:** None

---

### STORY-002: Implement Job Status Enum
**Story Points:** 2
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need a type-safe job status enum so that invalid states cannot be assigned.

**Acceptance Criteria:**
- [ ] Status is a string type alias
- [ ] Constants defined for: Pending, Running, Completed, Failed, Canceled
- [ ] String() method for display
- [ ] IsValid() method for validation
- [ ] IsTerminal() method to check if status is final

**Technical Tasks:**
- [ ] Create JobStatus type
- [ ] Define status constants
- [ ] Implement validation methods
- [ ] Write unit tests

**Dependencies:** None

---

### STORY-003: Implement State Machine Logic
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need state transition validation so that jobs follow valid lifecycle paths.

**Acceptance Criteria:**
- [ ] Cannot transition from terminal states (completed, failed, canceled)
- [ ] Pending → Running allowed
- [ ] Running → Completed allowed
- [ ] Running → Failed allowed
- [ ] Running → Canceled allowed
- [ ] Pending → Canceled allowed
- [ ] All other transitions rejected with error
- [ ] State transitions are atomic

**Technical Tasks:**
- [ ] Create CanTransition(from, to Status) method
- [ ] Implement state transition validation logic
- [ ] Add TransitionTo(status Status) error method on Job
- [ ] Write comprehensive tests for all valid/invalid transitions
- [ ] Document state machine with diagram

**Dependencies:** STORY-002

---

### STORY-004: Add Job Business Logic Methods
**Story Points:** 3
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need methods to manipulate job state according to business rules.

**Acceptance Criteria:**
- [ ] IncrementAttempts() increases attempts counter
- [ ] SetError(err) captures last error message
- [ ] ShouldRetry() returns true if attempts < MaxAttempts
- [ ] IsRetryable() checks if job can be retried
- [ ] Methods enforce business rules

**Technical Tasks:**
- [ ] Implement IncrementAttempts()
- [ ] Implement SetError(error)
- [ ] Implement ShouldRetry() bool
- [ ] Implement IsRetryable() bool
- [ ] Write unit tests for each method

**Dependencies:** STORY-001, STORY-003

---

## Epic 2: Repository Layer

### STORY-005: Define Repository Interface
**Story Points:** 2
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need a repository interface so that job storage can be implemented with different backends.

**Acceptance Criteria:**
- [ ] Interface includes: Create, GetByID, Update, Delete, List
- [ ] List supports filtering by status
- [ ] All methods return appropriate errors
- [ ] Interface is storage-agnostic

**Technical Tasks:**
- [ ] Create `domain/repository.go`
- [ ] Define Repository interface
- [ ] Define query options struct
- [ ] Document interface methods
- [ ] Define custom error types (ErrNotFound, etc.)

**Dependencies:** STORY-001

---

### STORY-006: Implement In-Memory Repository
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need a thread-safe in-memory repository so that jobs can be stored and retrieved during development.

**Acceptance Criteria:**
- [ ] Implements Repository interface
- [ ] Thread-safe using sync.RWMutex
- [ ] Support Create, GetByID, Update, Delete, List operations
- [ ] List can filter by status
- [ ] Returns ErrNotFound when job doesn't exist

**Technical Tasks:**
- [ ] Create `infrastructure/memory/repository.go`
- [ ] Implement NewInMemoryRepository()
- [ ] Implement all Repository interface methods
- [ ] Use sync.RWMutex for concurrency control
- [ ] Write unit tests with concurrent access
- [ ] Run tests with `-race` flag

**Dependencies:** STORY-005

---

### STORY-007: Add Repository Query Capabilities
**Story Points:** 3
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need to query jobs by various criteria so that the dispatcher can find pending jobs.

**Acceptance Criteria:**
- [ ] ListByStatus(status) returns jobs with given status
- [ ] ListPending() returns all pending jobs
- [ ] Results sorted by CreatedAt
- [ ] Empty slice returned if no matches

**Technical Tasks:**
- [ ] Add ListByStatus method to interface and implementation
- [ ] Add ListPending convenience method
- [ ] Implement sorting logic
- [ ] Write unit tests for queries

**Dependencies:** STORY-006

---

## Epic 3: Worker Pool & Concurrency

### STORY-008: Create Worker Pool Structure
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a developer, I need a worker pool that manages multiple goroutines so that jobs can be executed concurrently.

**Acceptance Criteria:**
- [ ] WorkerPool struct with configurable size
- [ ] Start() method spawns worker goroutines
- [ ] Stop() method gracefully shuts down workers
- [ ] Uses sync.WaitGroup to track workers
- [ ] Job channel for distributing work

**Technical Tasks:**
- [ ] Create `internal/worker/pool.go`
- [ ] Define WorkerPool struct
- [ ] Implement NewWorkerPool(size int)
- [ ] Implement Start() and Stop() methods
- [ ] Add job channel (buffered)
- [ ] Write tests for start/stop lifecycle

**Dependencies:** None

---

### STORY-009: Implement Worker Execution Loop
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a worker, I need to continuously process jobs from the queue so that jobs are executed.

**Acceptance Criteria:**
- [ ] Worker receives jobs from channel
- [ ] Worker executes job handler
- [ ] Worker respects context cancellation
- [ ] Worker handles panics gracefully
- [ ] Worker logs execution events

**Technical Tasks:**
- [ ] Implement worker() goroutine function
- [ ] Use select for job channel and context.Done()
- [ ] Add panic recovery with defer/recover
- [ ] Log worker start, job execution, errors
- [ ] Write tests for worker execution

**Dependencies:** STORY-008

---

### STORY-010: Create Job Dispatcher
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a system, I need a dispatcher that fetches pending jobs and sends them to workers so that jobs are processed.

**Acceptance Criteria:**
- [ ] Dispatcher polls repository for pending jobs
- [ ] Jobs sent to worker pool channel
- [ ] Polling interval configurable
- [ ] Dispatcher respects context cancellation
- [ ] Handles case when channel is full (backpressure)

**Technical Tasks:**
- [ ] Create `internal/dispatcher/dispatcher.go`
- [ ] Implement Dispatcher struct with repository dependency
- [ ] Implement Start() method with polling loop
- [ ] Add configurable polling interval
- [ ] Handle channel full scenario
- [ ] Write integration tests

**Dependencies:** STORY-007, STORY-008

---

### STORY-011: Add Worker Pool Configuration
**Story Points:** 2
**Priority:** P0
**Type:** Task

**Description:**
As an operator, I need to configure the worker pool size so that I can optimize resource usage.

**Acceptance Criteria:**
- [ ] Worker pool size configurable
- [ ] Channel buffer size configurable
- [ ] Default values provided
- [ ] Validation of configuration (size > 0)

**Technical Tasks:**
- [ ] Create Config struct for worker pool
- [ ] Add validation function
- [ ] Use functional options pattern
- [ ] Document configuration options

**Dependencies:** STORY-008

---

### STORY-012: Test Concurrent Job Execution
**Story Points:** 3
**Priority:** P0
**Type:** Test

**Description:**
As a developer, I need to verify that concurrent job execution works correctly without race conditions.

**Acceptance Criteria:**
- [ ] Multiple jobs execute concurrently
- [ ] No race conditions (verified with -race flag)
- [ ] Jobs complete successfully
- [ ] Worker pool handles load correctly

**Technical Tasks:**
- [ ] Write integration test with 100+ jobs
- [ ] Run tests with `-race` flag
- [ ] Verify all jobs complete
- [ ] Test with various pool sizes

**Dependencies:** STORY-008, STORY-009, STORY-010

---

## Epic 4: Job Execution Flow

### STORY-013: Implement Job Status Update Service
**Story Points:** 3
**Priority:** P0
**Type:** Task

**Description:**
As a worker, I need a service to update job status atomically so that status changes are consistent.

**Acceptance Criteria:**
- [ ] UpdateStatus(jobID, status) method
- [ ] Validates state transition before update
- [ ] Updates UpdatedAt timestamp
- [ ] Returns error if transition invalid
- [ ] Thread-safe operation

**Technical Tasks:**
- [ ] Create `internal/service/job_service.go`
- [ ] Implement JobService struct
- [ ] Add UpdateStatus method with validation
- [ ] Use repository for persistence
- [ ] Write unit tests

**Dependencies:** STORY-003, STORY-006

---

### STORY-014: Implement Job Execution Workflow
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a worker, I need to execute the complete job workflow from pending to completion so that jobs are processed correctly.

**Acceptance Criteria:**
- [ ] Fetch job from repository
- [ ] Update status to Running
- [ ] Execute job handler
- [ ] On success: update to Completed
- [ ] On failure: update to Failed with error
- [ ] Update UpdatedAt timestamp

**Technical Tasks:**
- [ ] Create ExecuteJob(ctx, jobID) method
- [ ] Implement status updates at each step
- [ ] Add error handling
- [ ] Update timestamps
- [ ] Write integration tests

**Dependencies:** STORY-013

---

### STORY-015: Add Panic Recovery in Job Execution
**Story Points:** 2
**Priority:** P0
**Type:** Task

**Description:**
As a system, I need to recover from panics during job execution so that one failing job doesn't crash the worker.

**Acceptance Criteria:**
- [ ] Panic during job execution is recovered
- [ ] Job marked as Failed
- [ ] Panic message captured in LastError
- [ ] Worker continues processing other jobs
- [ ] Panic logged

**Technical Tasks:**
- [ ] Add defer/recover in job execution
- [ ] Convert panic to error
- [ ] Update job status to Failed
- [ ] Log panic with stack trace
- [ ] Write test that triggers panic

**Dependencies:** STORY-014

---

### STORY-016: Integrate Worker Pool with Job Execution
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a system, I need workers to execute the full job workflow so that jobs move through their lifecycle.

**Acceptance Criteria:**
- [ ] Workers receive job IDs from channel
- [ ] Workers execute full job workflow
- [ ] Multiple workers process jobs concurrently
- [ ] Job status correctly updated in repository

**Technical Tasks:**
- [ ] Modify worker loop to call ExecuteJob
- [ ] Pass context to job execution
- [ ] Handle execution errors
- [ ] Write end-to-end integration test

**Dependencies:** STORY-009, STORY-014

---

## Epic 5: Retry Mechanism

### STORY-017: Implement Backoff Strategy
**Story Points:** 3
**Priority:** P0
**Type:** Task

**Description:**
As a system, I need a backoff strategy for retries so that failed jobs don't overwhelm downstream services.

**Acceptance Criteria:**
- [ ] Exponential backoff algorithm implemented
- [ ] Configurable base delay and max delay
- [ ] Optional jitter to prevent thundering herd
- [ ] CalculateBackoff(attempt int) time.Duration method

**Technical Tasks:**
- [ ] Create `internal/retry/backoff.go`
- [ ] Implement exponential backoff function
- [ ] Add jitter calculation
- [ ] Make configurable
- [ ] Write unit tests

**Dependencies:** None

---

### STORY-018: Add Retry Logic to Job Execution
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a system, I need to automatically retry failed jobs so that transient failures are handled gracefully.

**Acceptance Criteria:**
- [ ] When job fails, check if ShouldRetry()
- [ ] If yes: increment attempts, schedule retry
- [ ] If no: mark as terminal Failed
- [ ] Use backoff strategy for retry delay
- [ ] Capture error in LastError

**Technical Tasks:**
- [ ] Modify job execution to check retry eligibility
- [ ] Implement retry scheduling
- [ ] Use time.AfterFunc or ticker for delay
- [ ] Update attempts counter
- [ ] Write integration test with failing job

**Dependencies:** STORY-014, STORY-017

---

### STORY-019: Create Retry Scheduler
**Story Points:** 5
**Priority:** P0
**Type:** Task

**Description:**
As a system, I need a scheduler to re-queue jobs for retry after backoff period so that retries happen automatically.

**Acceptance Criteria:**
- [ ] Schedule retry after calculated backoff
- [ ] Reset job status to Pending
- [ ] Add job back to queue
- [ ] Scheduler respects system shutdown

**Technical Tasks:**
- [ ] Create retry scheduler component
- [ ] Implement ScheduleRetry(job, delay)
- [ ] Reset status to Pending before re-queueing
- [ ] Handle cancellation during backoff
- [ ] Write tests for retry scheduling

**Dependencies:** STORY-017, STORY-018

---

### STORY-020: Test Retry Scenarios
**Story Points:** 3
**Priority:** P0
**Type:** Test

**Description:**
As a developer, I need to verify retry logic works correctly under various failure scenarios.

**Acceptance Criteria:**
- [ ] Job retries correct number of times
- [ ] Backoff delays are correct
- [ ] Final failure marked correctly
- [ ] Transient failures eventually succeed

**Technical Tasks:**
- [ ] Test job that fails then succeeds
- [ ] Test job that fails MaxAttempts times
- [ ] Verify backoff timing
- [ ] Test with different MaxAttempts values

**Dependencies:** STORY-018, STORY-019

---

## Epic 6: Job Cancellation

### STORY-021: Implement Job Cancellation API
**Story Points:** 3
**Priority:** P1
**Type:** Task

**Description:**
As an operator, I need an API to cancel jobs so that unwanted jobs can be stopped.

**Acceptance Criteria:**
- [ ] CancelJob(jobID) method
- [ ] Cannot cancel jobs in terminal states
- [ ] Pending jobs marked as Canceled immediately
- [ ] Running jobs signaled for cancellation
- [ ] Returns error if job not found or not cancelable

**Technical Tasks:**
- [ ] Add CancelJob method to JobService
- [ ] Validate job is cancelable
- [ ] Update status to Canceled for pending jobs
- [ ] Signal running jobs (next story)
- [ ] Write unit tests

**Dependencies:** STORY-013

---

### STORY-022: Context-Based Cancellation for Running Jobs
**Story Points:** 5
**Priority:** P1
**Type:** Task

**Description:**
As a worker, I need to respect context cancellation so that running jobs can be stopped.

**Acceptance Criteria:**
- [ ] Each job execution uses cancellable context
- [ ] Contexts stored in thread-safe map
- [ ] Cancel() called on context when job canceled
- [ ] Workers check context.Done()
- [ ] Job marked as Canceled when context canceled

**Technical Tasks:**
- [ ] Create context manager for running jobs
- [ ] Store contexts in sync.Map
- [ ] Modify ExecuteJob to use cancellable context
- [ ] Implement CancelRunningJob(jobID)
- [ ] Handle context cancellation in workers
- [ ] Write tests for running job cancellation

**Dependencies:** STORY-021

---

### STORY-023: Test Cancellation Scenarios
**Story Points:** 2
**Priority:** P1
**Type:** Test

**Description:**
As a developer, I need to verify cancellation works correctly for both pending and running jobs.

**Acceptance Criteria:**
- [ ] Pending job cancels immediately
- [ ] Running job stops when context canceled
- [ ] Cannot cancel completed jobs
- [ ] Cancellation idempotent

**Technical Tasks:**
- [ ] Test pending job cancellation
- [ ] Test running job cancellation (with long-running handler)
- [ ] Test canceling already completed job
- [ ] Test canceling same job twice

**Dependencies:** STORY-021, STORY-022

---

## Epic 7: Job Handler System

### STORY-024: Define Job Handler Interface
**Story Points:** 2
**Priority:** P1
**Type:** Task

**Description:**
As a developer, I need a handler interface so that different job types can have custom execution logic.

**Acceptance Criteria:**
- [ ] Handler interface with Execute method
- [ ] Execute receives context and payload
- [ ] Execute returns error on failure
- [ ] Interface documented

**Technical Tasks:**
- [ ] Create `internal/handler/handler.go`
- [ ] Define Handler interface
- [ ] Add Execute(ctx context.Context, payload []byte) error
- [ ] Document interface contract

**Dependencies:** None

---

### STORY-025: Implement Handler Registry
**Story Points:** 3
**Priority:** P1
**Type:** Task

**Description:**
As a system, I need a registry to map job types to handlers so that jobs execute the correct logic.

**Acceptance Criteria:**
- [ ] Registry with Register(jobType, handler) method
- [ ] Get(jobType) returns handler or error
- [ ] Thread-safe access
- [ ] Error if handler not found
- [ ] Error if duplicate registration

**Technical Tasks:**
- [ ] Create Registry struct
- [ ] Implement Register method
- [ ] Implement Get method
- [ ] Use sync.RWMutex for thread safety
- [ ] Write unit tests

**Dependencies:** STORY-024

---

### STORY-026: Create Sample Job Handlers
**Story Points:** 5
**Priority:** P1
**Type:** Task

**Description:**
As a developer, I need example handlers to demonstrate the system's capabilities.

**Acceptance Criteria:**
- [ ] SendEmailHandler (simulated with sleep and log)
- [ ] GenerateReportHandler (simulated computation)
- [ ] ProcessFileHandler (simulated file processing)
- [ ] Each handler respects context cancellation
- [ ] Handlers parse payload correctly

**Technical Tasks:**
- [ ] Create SendEmailHandler with payload struct
- [ ] Create GenerateReportHandler
- [ ] Create ProcessFileHandler
- [ ] Make handlers cancellable
- [ ] Write tests for each handler

**Dependencies:** STORY-024

---

### STORY-027: Integrate Handlers with Job Execution
**Story Points:** 3
**Priority:** P1
**Type:** Task

**Description:**
As a worker, I need to execute jobs using the appropriate handler so that job type determines execution logic.

**Acceptance Criteria:**
- [ ] Job execution looks up handler by job type
- [ ] Handler executed with job payload
- [ ] Error returned if handler not found
- [ ] Handler execution timed/logged

**Technical Tasks:**
- [ ] Modify ExecuteJob to use handler registry
- [ ] Look up handler by job.Type
- [ ] Call handler.Execute with context and payload
- [ ] Handle unknown job types gracefully
- [ ] Write integration test

**Dependencies:** STORY-025, STORY-026

---

## Epic 8: HTTP API

### STORY-028: Setup HTTP Server
**Story Points:** 2
**Priority:** P1
**Type:** Task

**Description:**
As a developer, I need an HTTP server so that the API can accept requests.

**Acceptance Criteria:**
- [ ] HTTP server with configurable port
- [ ] Graceful startup and shutdown
- [ ] Request logging middleware
- [ ] Error handling middleware

**Technical Tasks:**
- [ ] Create `internal/api/server.go`
- [ ] Setup HTTP server with net/http or gin
- [ ] Add logging middleware
- [ ] Add error handling middleware
- [ ] Write tests

**Dependencies:** None

---

### STORY-029: Implement Create Job Endpoint
**Story Points:** 5
**Priority:** P1
**Type:** Task

**Description:**
As a client, I need an endpoint to create jobs so that I can submit work to the system.

**Acceptance Criteria:**
- [ ] POST /jobs endpoint
- [ ] Request body: { type, payload, maxAttempts }
- [ ] Validates request body
- [ ] Creates job in repository
- [ ] Returns 201 with job details
- [ ] Returns 400 for validation errors

**Technical Tasks:**
- [ ] Create CreateJobRequest DTO
- [ ] Create CreateJobResponse DTO
- [ ] Implement POST /jobs handler
- [ ] Add request validation
- [ ] Write API tests

**Dependencies:** STORY-028

---

### STORY-030: Implement Get Job Endpoint
**Story Points:** 3
**Priority:** P1
**Type:** Task

**Description:**
As a client, I need an endpoint to retrieve job status so that I can track job progress.

**Acceptance Criteria:**
- [ ] GET /jobs/:id endpoint
- [ ] Returns job details with status
- [ ] Returns 404 if job not found
- [ ] Returns 200 with job JSON

**Technical Tasks:**
- [ ] Create JobResponse DTO
- [ ] Implement GET /jobs/:id handler
- [ ] Map Job entity to response
- [ ] Write API tests

**Dependencies:** STORY-028

---

### STORY-031: Implement Cancel Job Endpoint
**Story Points:** 3
**Priority:** P1
**Type:** Task

**Description:**
As a client, I need an endpoint to cancel jobs so that I can stop unwanted work.

**Acceptance Criteria:**
- [ ] DELETE /jobs/:id endpoint
- [ ] Cancels job if cancelable
- [ ] Returns 200 on success
- [ ] Returns 404 if not found
- [ ] Returns 409 if not cancelable

**Technical Tasks:**
- [ ] Implement DELETE /jobs/:id handler
- [ ] Call CancelJob service method
- [ ] Handle errors appropriately
- [ ] Write API tests

**Dependencies:** STORY-028, STORY-021

---

### STORY-032: Implement List Jobs Endpoint
**Story Points:** 5
**Priority:** P1
**Type:** Task

**Description:**
As a client, I need an endpoint to list jobs so that I can view all jobs in the system.

**Acceptance Criteria:**
- [ ] GET /jobs endpoint
- [ ] Support ?status= query parameter
- [ ] Support pagination with ?limit= and ?offset=
- [ ] Returns array of jobs
- [ ] Returns 200 with jobs array

**Technical Tasks:**
- [ ] Create ListJobsRequest with query params
- [ ] Create ListJobsResponse DTO
- [ ] Implement GET /jobs handler
- [ ] Add filtering and pagination logic
- [ ] Write API tests

**Dependencies:** STORY-028

---

### STORY-033: Add API Request Validation
**Story Points:** 3
**Priority:** P1
**Type:** Task

**Description:**
As a system, I need to validate API requests so that invalid data is rejected early.

**Acceptance Criteria:**
- [ ] Validate required fields
- [ ] Validate field types and formats
- [ ] Return 400 with clear error messages
- [ ] Use validation library (e.g., validator)

**Technical Tasks:**
- [ ] Add validation tags to DTOs
- [ ] Implement validation middleware
- [ ] Return structured error responses
- [ ] Write validation tests

**Dependencies:** STORY-029, STORY-030, STORY-031, STORY-032

---

### STORY-034: Create API Documentation
**Story Points:** 3
**Priority:** P1
**Type:** Documentation

**Description:**
As a client developer, I need API documentation so that I can understand how to use the API.

**Acceptance Criteria:**
- [ ] OpenAPI 3.0 spec created
- [ ] All endpoints documented
- [ ] Request/response schemas defined
- [ ] Example requests and responses
- [ ] Swagger UI available (optional)

**Technical Tasks:**
- [ ] Create openapi.yaml
- [ ] Document all endpoints
- [ ] Add request/response examples
- [ ] Setup Swagger UI (optional)

**Dependencies:** STORY-029, STORY-030, STORY-031, STORY-032

---

## Epic 9: Observability

### STORY-035: Implement Structured Logging
**Story Points:** 5
**Priority:** P1
**Type:** Task

**Description:**
As an operator, I need structured logs so that I can debug issues and monitor the system.

**Acceptance Criteria:**
- [ ] Use slog or zerolog
- [ ] JSON log format
- [ ] Include: timestamp, level, message, jobID, workerID
- [ ] Log job lifecycle events
- [ ] Configurable log level

**Technical Tasks:**
- [ ] Setup logging library
- [ ] Create logger wrapper/helper
- [ ] Add logging to all critical paths
- [ ] Include correlation IDs
- [ ] Configure log levels via env var

**Dependencies:** None

---

### STORY-036: Add Metrics Collection
**Story Points:** 5
**Priority:** P1
**Type:** Task

**Description:**
As an operator, I need metrics so that I can monitor system health and performance.

**Acceptance Criteria:**
- [ ] Metrics for jobs_created, jobs_completed, jobs_failed
- [ ] Metric for jobs_retried
- [ ] Gauge for active_workers
- [ ] Histogram for job_duration
- [ ] Prometheus format (optional)

**Technical Tasks:**
- [ ] Setup metrics library (expvar or prometheus)
- [ ] Define metric collectors
- [ ] Instrument job lifecycle
- [ ] Expose /metrics endpoint
- [ ] Write tests

**Dependencies:** None

---

### STORY-037: Create Health Check Endpoint
**Story Points:** 2
**Priority:** P1
**Type:** Task

**Description:**
As an operator, I need a health check endpoint so that I can verify the system is running.

**Acceptance Criteria:**
- [ ] GET /health endpoint
- [ ] Returns 200 if healthy
- [ ] Checks: workers running, dispatcher running
- [ ] Returns JSON with status details

**Technical Tasks:**
- [ ] Implement /health endpoint
- [ ] Check worker pool status
- [ ] Check dispatcher status
- [ ] Return structured response

**Dependencies:** STORY-028

---

## Epic 10: Database Persistence

### STORY-038: Create Database Schema
**Story Points:** 3
**Priority:** P2
**Type:** Task

**Description:**
As a developer, I need a database schema so that jobs can be stored persistently.

**Acceptance Criteria:**
- [ ] Jobs table with all job attributes
- [ ] Indexes on status and created_at
- [ ] Primary key on id (UUID)
- [ ] Appropriate column types

**Technical Tasks:**
- [ ] Design jobs table schema
- [ ] Create migration file
- [ ] Add indexes
- [ ] Document schema

**Dependencies:** None

---

### STORY-039: Setup Database Migrations
**Story Points:** 2
**Priority:** P2
**Type:** Task

**Description:**
As a developer, I need database migrations so that schema changes are versioned and automated.

**Acceptance Criteria:**
- [ ] Use golang-migrate or similar
- [ ] Create initial migration
- [ ] Migrations run on startup
- [ ] Support up and down migrations

**Technical Tasks:**
- [ ] Setup golang-migrate
- [ ] Create initial migration
- [ ] Add migration runner to app startup
- [ ] Write migration tests

**Dependencies:** STORY-038

---

### STORY-040: Implement PostgreSQL Repository
**Story Points:** 8
**Priority:** P2
**Type:** Task

**Description:**
As a developer, I need a PostgreSQL repository so that jobs are stored in a database.

**Acceptance Criteria:**
- [ ] Implements Repository interface
- [ ] Connection pooling configured
- [ ] All CRUD operations implemented
- [ ] Uses transactions where appropriate
- [ ] Proper error handling

**Technical Tasks:**
- [ ] Create `infrastructure/postgres/repository.go`
- [ ] Setup connection pool with pgx
- [ ] Implement all Repository methods
- [ ] Use parameterized queries
- [ ] Write integration tests with test database

**Dependencies:** STORY-038, STORY-039

---

### STORY-041: Add Job Claiming with SKIP LOCKED
**Story Points:** 5
**Priority:** P2
**Type:** Task

**Description:**
As a dispatcher, I need to claim pending jobs atomically so that multiple instances don't process the same job.

**Acceptance Criteria:**
- [ ] Use SELECT FOR UPDATE SKIP LOCKED
- [ ] ClaimJobs(limit) returns pending jobs
- [ ] Jobs atomically marked as claimed/running
- [ ] Supports multiple dispatcher instances

**Technical Tasks:**
- [ ] Implement ClaimJobs method
- [ ] Use FOR UPDATE SKIP LOCKED
- [ ] Update status in same transaction
- [ ] Write concurrency tests

**Dependencies:** STORY-040

---

## Epic 11: Configuration Management

### STORY-042: Externalize Configuration
**Story Points:** 5
**Priority:** P2
**Type:** Task

**Description:**
As an operator, I need externalized configuration so that I can deploy across environments without code changes.

**Acceptance Criteria:**
- [ ] Worker pool size configurable
- [ ] Database connection string configurable
- [ ] Retry settings configurable
- [ ] HTTP port configurable
- [ ] Support env vars and config file
- [ ] Sensible defaults provided

**Technical Tasks:**
- [ ] Create Config struct
- [ ] Use viper or envconfig
- [ ] Load from env vars and file
- [ ] Validate configuration on startup
- [ ] Document all options

**Dependencies:** None

---

## Epic 12: Graceful Shutdown

### STORY-043: Implement Graceful Shutdown
**Story Points:** 8
**Priority:** P2
**Type:** Task

**Description:**
As an operator, I need graceful shutdown so that in-flight jobs complete before the system stops.

**Acceptance Criteria:**
- [ ] Responds to SIGTERM/SIGINT
- [ ] Stops accepting new jobs
- [ ] Waits for in-flight jobs (with timeout)
- [ ] Closes database connections
- [ ] Logs shutdown events

**Technical Tasks:**
- [ ] Setup signal handling
- [ ] Implement shutdown coordinator
- [ ] Stop dispatcher and workers gracefully
- [ ] Use context cancellation
- [ ] Add shutdown timeout
- [ ] Write shutdown tests

**Dependencies:** STORY-008, STORY-010

---

## Epic 13: Production Deployment

### STORY-044: Create Dockerfile
**Story Points:** 3
**Priority:** P2
**Type:** Task

**Description:**
As an operator, I need a Docker image so that I can deploy the application in containers.

**Acceptance Criteria:**
- [ ] Multi-stage Dockerfile
- [ ] Optimized image size (<50MB)
- [ ] Uses distroless or alpine
- [ ] Health check configured
- [ ] Non-root user

**Technical Tasks:**
- [ ] Create Dockerfile
- [ ] Use multi-stage build
- [ ] Configure health check
- [ ] Run as non-root
- [ ] Test image build and run

**Dependencies:** None

---

### STORY-045: Create Docker Compose Setup
**Story Points:** 3
**Priority:** P2
**Type:** Task

**Description:**
As a developer, I need a Docker Compose setup so that I can run the system locally easily.

**Acceptance Criteria:**
- [ ] docker-compose.yml with app and postgres
- [ ] Environment variables configured
- [ ] Volume mounts for development
- [ ] Health checks
- [ ] Easy startup (docker-compose up)

**Technical Tasks:**
- [ ] Create docker-compose.yml
- [ ] Configure services (app, postgres)
- [ ] Setup volumes
- [ ] Add health checks
- [ ] Document usage

**Dependencies:** STORY-044

---

## Epic 14: Documentation & Examples

### STORY-046: Create Comprehensive README
**Story Points:** 5
**Priority:** P2
**Type:** Documentation

**Description:**
As a user, I need clear documentation so that I can understand and use the job processor.

**Acceptance Criteria:**
- [ ] Project overview and purpose
- [ ] Quick start guide
- [ ] Architecture overview
- [ ] Configuration options documented
- [ ] Examples included

**Technical Tasks:**
- [ ] Write project overview
- [ ] Create quick start guide
- [ ] Document architecture
- [ ] Add configuration reference
- [ ] Include examples

**Dependencies:** All previous stories

---

### STORY-047: Create Handler Development Guide
**Story Points:** 3
**Priority:** P2
**Type:** Documentation

**Description:**
As a developer, I need a guide for creating custom handlers so that I can extend the system.

**Acceptance Criteria:**
- [ ] Handler interface explained
- [ ] Step-by-step guide for creating handler
- [ ] Example handler code
- [ ] Testing guide

**Technical Tasks:**
- [ ] Write handler development guide
- [ ] Create example handler
- [ ] Document testing approach
- [ ] Add to documentation

**Dependencies:** STORY-024, STORY-025, STORY-026

---

### STORY-048: Create Operations Runbook
**Story Points:** 3
**Priority:** P2
**Type:** Documentation

**Description:**
As an operator, I need a runbook so that I can deploy, monitor, and troubleshoot the system.

**Acceptance Criteria:**
- [ ] Deployment procedures
- [ ] Monitoring guidelines
- [ ] Common troubleshooting scenarios
- [ ] Performance tuning tips

**Technical Tasks:**
- [ ] Document deployment process
- [ ] Create monitoring guide
- [ ] Document troubleshooting steps
- [ ] Add performance tuning section

**Dependencies:** STORY-042, STORY-043, STORY-044, STORY-045

---

## Summary

**Total Stories:** 48
**Total Story Points:** ~172

### By Epic
- Epic 1: 13 points (4 stories)
- Epic 2: 10 points (3 stories)
- Epic 3: 20 points (5 stories)
- Epic 4: 15 points (4 stories)
- Epic 5: 16 points (4 stories)
- Epic 6: 10 points (3 stories)
- Epic 7: 13 points (4 stories)
- Epic 8: 21 points (7 stories)
- Epic 9: 12 points (3 stories)
- Epic 10: 18 points (4 stories)
- Epic 11: 5 points (1 story)
- Epic 12: 8 points (1 story)
- Epic 13: 6 points (2 stories)
- Epic 14: 11 points (3 stories)

### MVP (Epics 1-8)
**31 stories, ~102 story points**
