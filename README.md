# Go Job Processor

## Idea

Build a lightweight **asynchronous job processing system** in Go.

The system receives jobs, executes them in the background using a worker pool, retries failures, supports cancellation, and keeps track of job state.

This project focuses on:

- Go concurrency (goroutines, channels, worker pools)
- Clean Architecture principles
- Failure handling and retries
- Clear domain modeling
- Production-like backend thinking

No UI, no auth, no fluff.

---

## Problem Statement

Many backend systems need to run tasks asynchronously:

- sending emails
- generating reports
- processing files
- syncing data

This project implements a **small but realistic job processor** similar in spirit to Sidekiq or Temporal (simplified).

---

## Core Capabilities (MVP)

- Create a job
- Execute jobs asynchronously
- Track job status
- Retry failed jobs
- Cancel a job
- Query job status

If these work, the project is considered **complete**.

---

## Job Domain

### Job Attributes

- ID
- Type
- Payload
- Status
- Attempts
- MaxAttempts
- LastError
- CreatedAt
- UpdatedAt

### Job Statuses

- `pending`
- `running`
- `completed`
- `failed`
- `canceled`

### State Rules

- Jobs start as `pending`
- Only one worker can execute a job at a time
- `completed`, `failed`, and `canceled` are terminal states
- Attempts increase only when a job fails
- A job cannot transition from a terminal state

---

## Architecture Overview

- **Domain**

  - Job entity
  - State transition rules

- **Application**

  - Use cases (create job, execute job, cancel job)
  - Worker pool and dispatcher
  - Concurrency logic

- **Infrastructure**
  - Job repository (start in memory, later database)
  - HTTP API (optional for MVP)

---

## Concurrency Model

- Dispatcher fetches pending jobs
- Jobs are sent to a buffered channel
- A fixed number of workers consume jobs
- Each worker:
  - Marks job as running
  - Executes the job handler
  - Marks job as completed or failed

Channels provide:

- Backpressure
- Simple coordination
- Clear ownership of work

---

## Job Execution

Jobs are executed via a handler interface:

- Each job type has its own handler
- Handlers receive a context and payload
- Handlers return an error on failure

This allows:

- Easy testing
- Clear separation of concerns
- Extensibility

---

## Retry & Cancellation

- Failed jobs are retried up to `MaxAttempts`
- Retry logic uses simple backoff
- Jobs can be canceled while pending or running
- Workers respect context cancellation

---

## Scope Control

This project intentionally avoids:

- Distributed queues
- Message brokers
- Authentication
- UI
- Kubernetes

The goal is **depth, not breadth**.

---

## First Milestone

**Goal:** Jobs move through the system.

- In-memory repository
- One simple job handler (e.g. sleep or print)
- Worker pool running jobs
- Status transitions visible via logs

Once this works, the foundation is complete.

---

## Why This Project Exists

- Strengthen Go concurrency skills
- Demonstrate senior-level backend design
- Create a strong portfolio artifact
- Have something concrete to explain in interviews

This is not a tutorial project.
It is a **learning-through-building system**.
