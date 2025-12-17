# Milestone 1: Foundation & Core Database

**Duration:** 2 weeks (Sprint 1)  
**Status:** Not Started  

---

## Milestone Objectives

1. Establish project structure and development workflow
2. Implement database layer with PostgreSQL
3. Create core job data operations (CRUD)
4. Set up testing and CI infrastructure
5. Build configuration system

---

## Success Criteria

- [ ] Project builds and tests run in CI
- [ ] Can create, read, update, delete jobs in database
- [ ] Configuration loads from environment variables and file
- [ ] Database migrations run successfully
- [ ] Test coverage > 80% on storage package

---

## Epic 1.1: Project Setup

**Points:** 5 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M1-001 | Initialize Go module with standard directory structure | 1 | `go build` succeeds; follows Go project layout conventions |
| M1-002 | Create Makefile with build, test, lint targets | 1 | `make build`, `make test`, `make lint` all work |
| M1-003 | Configure golangci-lint with team standards | 1 | Linting passes; `.golangci.yml` committed |
| M1-004 | Set up GitHub Actions for CI | 2 | Tests run on PR; lint check passes; build artifacts created |
| M1-005 | Add repository essentials (gitignore, license, contributing guide) | 0.5 | Files present and appropriate |

**Deliverable:** Repository ready for team development

---

## Epic 1.2: Database Schema & Migrations

**Points:** 8 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M1-006 | Evaluate and select migration tool (golang-migrate vs goose) | 1 | Decision documented; tool integrated |
| M1-007 | Create jobs table migration with all required columns | 3 | Migration applies cleanly; schema matches ADR-001 |
| M1-008 | Create indexes for job claiming and status queries | 2 | Partial index on pending jobs; status index exists |
| M1-009 | Add migration for schema versioning/metadata | 1 | Can track which migrations have run |
| M1-010 | Document migration workflow for developers | 0.5 | README in migrations folder; clear instructions |

**Schema Requirements:**
- Job ID (UUID, primary key)
- Queue name (string, indexed)
- Job type (string)
- Payload (JSON)
- Status (enum: pending, running, completed, dead, cancelled)
- Attempt tracking (current attempts, max attempts)
- Error tracking (last error message)
- Timestamps (created, updated, scheduled, started, completed)
- Idempotency key (unique, nullable)
- Worker lock tracking (locked_by, locked_at)

**Deliverable:** Database schema ready for job storage

---

## Epic 1.3: Storage Layer Implementation

**Points:** 13 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M1-011 | Define storage interface (contract for all backends) | 2 | Interface supports all required operations; clean abstraction |
| M1-012 | Implement PostgreSQL connection pool setup | 2 | Pool configurable; connections managed properly |
| M1-013 | Implement CreateJob operation | 2 | Jobs created with all fields; idempotency key enforced |
| M1-014 | Implement GetJob operation | 1 | Retrieves job by ID; returns not found appropriately |
| M1-015 | Implement ListJobs with filtering and pagination | 3 | Filter by status, queue, type; pagination works |
| M1-016 | Implement UpdateJob operation | 2 | Can update status, attempts, error; updated_at auto-set |
| M1-017 | Implement DeleteJob operation | 1 | Soft delete or hard delete; handles not found |

**Interface Methods Required:**
- `CreateJob(ctx, job) → error`
- `GetJob(ctx, id) → job, error`
- `GetJobByIdempotencyKey(ctx, key) → job, error`
- `ListJobs(ctx, filter) → jobs, error`
- `UpdateJob(ctx, id, updates) → error`
- `DeleteJob(ctx, id) → error`
- `Ping(ctx) → error`
- `Close() → error`

**Deliverable:** Working storage layer with PostgreSQL implementation

---

## Epic 1.4: Configuration System

**Points:** 5 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M1-018 | Design configuration structure (all settings identified) | 1 | Config struct covers server, database, worker, logging |
| M1-019 | Implement environment variable loading | 2 | All settings readable from env vars; follows 12-factor |
| M1-020 | Implement config file loading (YAML) | 1 | Can load from goqueue.yaml; env vars override file |
| M1-021 | Add configuration validation on startup | 1 | Server fails fast on invalid config; clear error messages |

**Configuration Sections:**
- Server: host, port, timeouts
- Database: connection URL, pool settings
- Worker: concurrency, queues, poll interval, timeouts
- Logging: level, format (json/text)

**Deliverable:** Flexible configuration system

---

## Epic 1.5: Testing Infrastructure

**Points:** 5 | **Priority:** P0

| Task ID | Task | Points | Acceptance Criteria |
|---------|------|--------|---------------------|
| M1-022 | Set up test database helpers (test containers or local) | 2 | Tests can spin up isolated Postgres; cleanup automatic |
| M1-023 | Write unit tests for storage layer | 2 | All CRUD operations tested; edge cases covered |
| M1-024 | Create integration test framework | 1 | Can run end-to-end tests with real database |

**Testing Approach:**
- Use testcontainers-go for isolated PostgreSQL per test
- Each test gets fresh database with migrations applied
- Cleanup happens automatically after test

**Deliverable:** Reliable test infrastructure

---

## Sprint Schedule

### Week 1

| Day | Focus | Tasks |
|-----|-------|-------|
| Mon | Project scaffolding | M1-001, M1-002 |
| Tue | Tooling setup | M1-003, M1-004, M1-005 |
| Wed | Migration tooling | M1-006 |
| Thu | Schema design | M1-007, M1-008 |
| Fri | Schema finalization | M1-009, M1-010 |

### Week 2

| Day | Focus | Tasks |
|-----|-------|-------|
| Mon | Storage interface | M1-011, M1-012 |
| Tue | Core CRUD | M1-013, M1-014 |
| Wed | List operations | M1-015, M1-016 |
| Thu | Config system | M1-017, M1-018, M1-019 |
| Fri | Testing & polish | M1-020, M1-021, M1-022, M1-023, M1-024 |

---

## Dependencies

**External:**
- PostgreSQL 14+ available for development
- GitHub repository created
- CI/CD minutes available

**Blockers:**
- ADR-001 (storage backend) must be approved before starting

---

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Migration tool doesn't meet needs | Medium | Timebox evaluation to 4 hours; have backup choice |
| Schema design needs iteration | Low | Get early review; keep migrations reversible |
| Test containers slow in CI | Low | Consider caching; parallel test execution |

---

## Handoff to M2

At milestone completion, M2 can begin with:
- Working storage layer to build worker on top of
- Test infrastructure to validate worker behavior
- Configuration system to add worker settings
