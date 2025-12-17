# Milestone 2: Job Processing Engine

**Duration:** 2 weeks (Sprint 2)  
**Status:** Not Started  
**Depends On:** M1 Complete  

---

## Milestone Objectives

1. Implement worker pool with configurable concurrency
2. Build reliable job claiming with PostgreSQL SKIP LOCKED
3. Add retry logic with exponential backoff
4. Implement graceful shutdown
5. Handle dead letters (exhausted retries)

---

## Success Criteria

- [ ] Jobs claim atomically (no double-processing under load)
- [ ] Configurable worker concurrency working
- [ ] Failed jobs retry with exponential backoff
- [ ] Graceful shutdown completes in-flight jobs
- [ ] Dead jobs queryable after max retries exhausted
- [ ] Zero job loss on graceful restart
- [ ] No race conditions (verified with `-race` flag)

---

## Epic 2.1: Job Claiming

**Points:** 8 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M2-001 | Implement ClaimJob with SKIP LOCKED | 3 | Atomic claim; no job processed twice under concurrent load |
| M2-002 | Add worker ID tracking (locked_by field) | 1 | Can identify which worker instance has each job |
| M2-003 | Implement claim timeout for stuck jobs | 2 | Jobs locked > timeout return to pending automatically |
| M2-004 | Add batch claim support (claim N jobs at once) | 2 | Reduces database round trips; improves throughput |

**Claiming Behavior:**
- Single atomic query claims job and updates status to 'running'
- SKIP LOCKED prevents blocking on contested jobs
- Jobs claimed in scheduled_at order (FIFO within queue)
- Worker ID stored for debugging and monitoring

**Deliverable:** Reliable, concurrent-safe job claiming

---

## Epic 2.2: Worker Pool

**Points:** 13 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M2-005 | Define Worker interface and configuration struct | 1 | Clean abstraction; all settings configurable |
| M2-006 | Implement worker pool manager (spawn/manage goroutines) | 3 | Can start N workers; clean lifecycle management |
| M2-007 | Implement dispatcher (polls DB, routes to workers) | 3 | Efficiently distributes jobs to available workers |
| M2-008 | Add handler registration system | 2 | Register function per job type; unknown types fail cleanly |
| M2-009 | Implement polling with exponential backoff | 2 | Fast when busy; backs off when idle; efficient DB usage |
| M2-010 | Add per-queue concurrency configuration | 2 | Different queues can have different worker counts |

**Architecture:**
- One dispatcher goroutine polls database
- Worker pool of N goroutines process jobs
- Jobs distributed via Go channels
- Backpressure when all workers busy

**Polling Behavior:**
- Base interval: 100ms when active
- Backoff to 5s when queue empty
- Immediate poll on successful claim

**Deliverable:** Scalable worker pool architecture

---

## Epic 2.3: Job Execution

**Points:** 8 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M2-011 | Implement job execution with configurable timeout | 2 | Jobs exceeding timeout are cancelled and failed |
| M2-012 | Add panic recovery in job handlers | 2 | Handler panic doesn't crash worker; job marked failed |
| M2-013 | Implement job completion flow | 1 | Successful jobs marked completed with timestamp |
| M2-014 | Implement job failure handling | 2 | Failed jobs capture error; trigger retry or dead logic |
| M2-015 | Add context cancellation support | 1 | Handlers receive cancellable context; respect shutdown |

**Execution Flow:**
1. Job received from dispatcher channel
2. Create context with timeout
3. Look up handler by job type
4. Execute handler with panic recovery
5. On success: mark completed
6. On failure: trigger retry logic
7. On panic: capture stack trace, mark failed

**Deliverable:** Robust job execution with safety guarantees

---

## Epic 2.4: Retry Logic

**Points:** 8 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M2-016 | Implement exponential backoff calculation | 2 | Formula matches ADR-004; configurable parameters |
| M2-017 | Add jitter to backoff (±25%) | 1 | Prevents thundering herd on mass failures |
| M2-018 | Implement retry scheduling (reschedule failed jobs) | 2 | Failed jobs get new scheduled_at with backoff |
| M2-019 | Add custom retry timing support (RetryAfter) | 2 | Handlers can specify exact retry delay |
| M2-020 | Implement permanent failure detection | 1 | Handlers can mark errors as non-retryable |

**Retry Behavior:**
- Default: exponential backoff with 2x multiplier
- Base interval: 15 seconds
- Max interval: 1 hour
- Jitter: ±25% randomization
- Max attempts: 10 (configurable per job)

**Error Types:**
- Normal error → retry with backoff
- RetryAfter error → retry after specified duration
- Permanent error → skip retries, go to dead

**Deliverable:** Intelligent retry system

---

## Epic 2.5: Graceful Shutdown

**Points:** 5 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M2-021 | Implement OS signal handling (SIGTERM, SIGINT) | 1 | Signals trigger shutdown sequence |
| M2-022 | Stop dispatcher from claiming new jobs | 1 | No new jobs claimed after shutdown initiated |
| M2-023 | Wait for in-flight jobs with timeout | 2 | Configurable timeout; logs progress |
| M2-024 | Return incomplete jobs to pending state | 1 | Jobs in channel returned to queue; no job loss |

**Shutdown Sequence:**
1. Receive shutdown signal
2. Stop dispatcher polling
3. Close job channel (no new jobs to workers)
4. Wait for workers to finish (with timeout)
5. If timeout: return remaining jobs to pending
6. Close database connections
7. Exit cleanly

**Deliverable:** Zero job loss on restarts

---

## Epic 2.6: Dead Letter Handling

**Points:** 3 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M2-025 | Implement transition to dead status | 1 | Jobs exceeding max_attempts marked dead |
| M2-026 | Add dead job query support in storage | 1 | Can list dead jobs by queue for investigation |
| M2-027 | Implement manual retry for dead jobs | 1 | API to retry dead job (resets attempts, back to pending) |

**Dead Job Behavior:**
- Status set to 'dead' (not deleted)
- Final error message preserved
- Completion timestamp recorded
- Remains queryable for debugging
- Can be manually retried via API

**Deliverable:** Dead letter queue functionality

---

## Sprint Schedule

### Week 3

| Day | Focus | Tasks |
|-----|-------|-------|
| Mon | Job claiming | M2-001, M2-002 |
| Tue | Worker pool setup | M2-005, M2-006 |
| Wed | Dispatcher | M2-007, M2-009 |
| Thu | Handler system | M2-008 |
| Fri | Job execution | M2-011, M2-012 |

### Week 4

| Day | Focus | Tasks |
|-----|-------|-------|
| Mon | Completion/failure flows | M2-013, M2-014, M2-015 |
| Tue | Retry logic | M2-016, M2-017, M2-018 |
| Wed | Advanced retry | M2-019, M2-020 |
| Thu | Graceful shutdown | M2-021, M2-022, M2-023, M2-024 |
| Fri | Dead letters & testing | M2-025, M2-026, M2-027, M2-003, M2-004 |

---

## Test Scenarios

| Scenario | Expected Behavior |
|----------|-------------------|
| Submit 1000 jobs concurrently | All processed exactly once |
| Kill worker during job execution | Job returns to queue, processes on restart |
| Handler panics | Job fails, worker continues processing other jobs |
| Graceful shutdown during processing | In-flight jobs complete or return to queue |
| Max retries exceeded | Job moves to dead status with error preserved |
| Custom RetryAfter returned | Next attempt scheduled at specified time |

---

## Performance Targets

| Metric | Target |
|--------|--------|
| Job throughput | > 1,000 jobs/sec with 10 workers |
| Claim latency | < 5ms p99 |
| Job loss on graceful shutdown | 0% |
| Job loss on hard kill | Return to queue on restart |

---

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Race conditions in claiming | High | Extensive testing with `-race`; load testing |
| Deadlocks in shutdown | High | Careful channel design; timeout all waits |
| Memory leaks from goroutines | Medium | Context cancellation everywhere; leak testing |

---

## Handoff to M3

At milestone completion, M3 can begin with:
- Complete job processing engine to expose via API
- Storage layer extended with claiming operations
- Configuration extended with worker settings
