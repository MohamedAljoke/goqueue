# ADR-001: Storage Backend Selection

**Status:** Accepted  
**Date:** 2024-01-XX  
**Deciders:** Tech Lead, Backend Team  

---

## Context

We need a storage backend for persisting jobs. Requirements:
- Atomic job claiming (prevent double-processing)
- Efficient polling for available jobs
- Persistence across restarts
- Query capabilities for job status and observability

---

## Options Considered

### Option 1: PostgreSQL ✅ SELECTED

Uses PostgreSQL with `SELECT ... FOR UPDATE SKIP LOCKED` for atomic job claiming.

| Pros | Cons |
|------|------|
| ACID guarantees | Slightly higher latency than Redis (~5ms vs ~1ms) |
| Teams usually already have it | Polling-based (not push) |
| Rich querying for observability | Connection pool management needed |
| `SKIP LOCKED` designed for queues | |
| No additional infrastructure | |
| Battle-tested at scale | |

**Expected Performance:** 5,000-10,000 jobs/sec with proper indexing

---

### Option 2: Redis

Uses Redis lists with BRPOPLPUSH or Redis Streams.

| Pros | Cons |
|------|------|
| Very fast (sub-ms latency) | Persistence requires careful configuration |
| Push-based (blocking pop) | Additional infrastructure to manage |
| Native queue data structures | Less queryable for observability |
| Industry standard for job queues | Memory-bound storage |

**Expected Performance:** 50,000+ jobs/sec

---

### Option 3: SQLite (Embedded)

Embedded SQLite with WAL mode for single-node deployments.

| Pros | Cons |
|------|------|
| Zero dependencies | Single-node only forever |
| Single file deployment | Write contention at scale |
| Fast for low volume | No remote workers possible |

**Expected Performance:** ~1,000 jobs/sec

---

### Option 4: Embedded KV Store (BoltDB/BadgerDB)

Pure Go embedded key-value store.

| Pros | Cons |
|------|------|
| Pure Go, no CGO | Custom indexing needed |
| Very fast writes | Less familiar to operators |
| Single binary possible | Single-node only |
| | No SQL for ad-hoc queries |

---

## Decision

**PostgreSQL** as the primary backend for MVP.

### Rationale

1. **Zero additional infrastructure** - most teams already have Postgres
2. **Purpose-built feature** - SKIP LOCKED was literally designed for job queues
3. **Queryability** - easy observability, debugging, ad-hoc analysis
4. **Reliability** - ACID guarantees mean we can promise job delivery
5. **Good enough performance** - 5k jobs/sec covers 99% of use cases
6. **Future-proof** - can add Redis backend later for high-throughput needs

### Trade-offs Accepted

- Higher latency than Redis (~5ms vs ~1ms) - acceptable for background jobs
- Polling instead of push - will implement efficient polling with backoff
- Connection pool complexity - will provide sensible defaults

---

## Consequences

- Must implement efficient polling with exponential backoff when queue is empty
- Must design clean storage interface to allow Redis backend later
- Must document PostgreSQL version requirements (14+ recommended)
- Must provide good connection pool defaults in configuration
- Schema design must optimize for the "claim pending job" query

---

## Future Considerations

- Add Redis backend post-MVP for teams needing >10k jobs/sec
- Consider SQLite option for single-node/edge deployments
- Evaluate CockroachDB for distributed PostgreSQL needs
