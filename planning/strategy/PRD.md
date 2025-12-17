# Product Requirements Document (PRD)
## GoQueue - Golang Job Processor

**Version:** 0.1.0-MVP  
**Status:** Discovery  
**Owner:** Product Owner  

---

## 1. Problem Statement

Development teams using Go need a simple, reliable way to process background jobs. Current options are either too complex (AWS Batch, Temporal), tied to other ecosystems (Sidekiq/Ruby, Celery/Python), or require additional infrastructure (Faktory, RabbitMQ).

**We believe** a Go-native, single-binary job processor will reduce operational overhead and improve developer experience for Go teams.

---

## 2. Vision

A lightweight, self-hosted job processing system that:
- Deploys as a single binary
- Uses PostgreSQL (which teams already have)
- Provides reliability without complexity
- Scales from MVP to production workloads

---

## 3. Target Users

### Primary: Backend Developer (Go teams)
- Building microservices in Go
- Needs background processing (emails, reports, webhooks, data pipelines)
- Wants minimal operational overhead
- Values: simplicity, reliability, good DX

### Secondary: DevOps/Platform Engineer
- Deploys and monitors the system
- Wants single binary, easy configuration
- Values: low maintenance, good metrics, easy scaling

---

## 4. Competitive Analysis

| Solution | Pros | Cons | Our Differentiation |
|----------|------|------|---------------------|
| AWS Batch | Managed, scalable | Vendor lock-in, complex, expensive | Self-hosted, simple |
| Sidekiq | Proven, feature-rich | Ruby only | Go-native |
| Celery | Mature, flexible | Python only, complex setup | Single binary |
| Temporal | Powerful workflows | Steep learning curve, heavy | Simple jobs first |
| Faktory | Language agnostic | External daemon required | Embedded library |

---

## 5. MVP Scope

### 5.1 In Scope (Must Have)

| Feature | Description | Why MVP |
|---------|-------------|---------|
| Job Submission | HTTP API to enqueue jobs | Core functionality |
| Job Processing | Workers pull and execute jobs | Core functionality |
| Multiple Queues | Separate workloads by queue name | Basic organization |
| Retry Logic | Automatic retries with exponential backoff | Reliability |
| Job Status | Query job state via API | Observability |
| Persistence | Jobs survive server restarts | Reliability |
| CLI Tool | Submit jobs, check status from terminal | Developer experience |
| Prometheus Metrics | Queue depth, throughput, errors | Production readiness |
| Graceful Shutdown | Complete in-flight jobs on shutdown | Reliability |

### 5.2 Out of Scope (Post-MVP)

| Feature | Reason for Deferral |
|---------|---------------------|
| Web UI Dashboard | CLI sufficient for MVP; adds frontend complexity |
| Cron/Scheduled Jobs | Can use external cron initially; adds scheduler complexity |
| Job Dependencies (DAG) | Workflow orchestration is a different product |
| Priority Queues | Multiple named queues sufficient initially |
| Rate Limiting | Can be implemented in worker handlers |
| Redis Backend | PostgreSQL covers 90% of use cases |
| gRPC API | HTTP simpler for adoption; gRPC can come later |
| Multi-tenancy | Enterprise feature |
| Distributed Tracing | Nice-to-have; can add OpenTelemetry later |

### 5.3 Design Decisions Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Storage Backend | PostgreSQL | Teams already have it; SKIP LOCKED is perfect for queues |
| API Protocol | REST/HTTP | Universal, debuggable, lower adoption barrier |
| Worker Model | Embedded Go library | Best DX for Go teams; no extra processes |
| Retry Strategy | Exponential backoff with jitter | Industry standard; prevents thundering herd |

---

## 6. User Stories (Summary)

### Job Submission
- Submit a job via HTTP API with custom payload
- Specify which queue a job belongs to
- Set maximum retry attempts per job
- Use idempotency keys to prevent duplicates

### Job Processing
- Register handler functions for job types
- Configure concurrency per worker
- Jobs that panic are caught and marked failed
- Context cancellation respected for graceful shutdown

### Reliability
- Failed jobs retry automatically with backoff
- Jobs persist across server restarts
- Dead jobs (exhausted retries) preserved for debugging
- Manual retry of dead jobs supported

### Observability
- Query job status via API
- Prometheus metrics for dashboards/alerts
- View error messages for failed jobs
- List jobs filtered by status/queue

### Operations
- Start server with single binary
- Configure via environment variables or config file
- Health check endpoints for load balancers
- Graceful shutdown on SIGTERM

---

## 7. Success Metrics

| Metric | Target | How to Measure |
|--------|--------|----------------|
| Job Throughput | > 1,000 jobs/sec | Benchmark suite |
| Submission Latency | < 10ms p99 | Prometheus histogram |
| Job Loss Rate | 0% on graceful shutdown | Integration tests |
| Time to Hello World | < 15 minutes | User testing |
| Binary Size | < 20MB | Build pipeline |

---

## 8. Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| PostgreSQL becomes bottleneck | High | Medium | Design pluggable backend interface from start |
| Worker SDK too complex | High | Medium | Prioritize DX; user testing; extensive examples |
| Feature creep delays MVP | High | High | Strict scope discipline; defer to backlog |
| Competitor releases similar | Medium | Low | Move fast; focus on simplicity differentiation |

---

## 9. Open Questions

- [ ] Should jobs have TTL (auto-expire)?
- [ ] Dead letter queue: separate table or status flag?
- [ ] Default max retries: 10 or configurable per queue?
- [ ] Job payload size limit?

---

## 10. Timeline

**Target MVP Duration:** 8 weeks with 2-3 engineers

| Phase | Duration | Outcome |
|-------|----------|---------|
| Foundation | 2 weeks | Database layer, project setup |
| Processing Engine | 2 weeks | Workers, retries, reliability |
| API & SDK | 2 weeks | HTTP API, Go client, CLI |
| Operations | 2 weeks | Metrics, docs, packaging |
