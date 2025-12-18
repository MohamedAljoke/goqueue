# GoQueue Development Journey - Learning Path

A step-by-step guide to building a job queue system in Go from scratch.

## Learning Approach

This project follows a TDD and incremental refactoring approach:
1. Start simple with hardcoded implementations
2. Write tests to prove functionality
3. Refactor and improve design
4. Extract into packages when complexity grows

## Steps

### Completed

1. [Basic Job Processing](./step1-basic-job-processing.md)
   - Create Job struct
   - Simple hardcoded handler
   - First tests

2. [Configurable Handlers](./step2-configurable-handlers.md)
   - HandlerFunc type
   - Dependency injection
   - Function types in Go

3. [Error Handling](./step3-error-handling.md)
   - Check handler errors
   - Failed status
   - Error wrapping

4. [Storage Layer](./step4-add-storage.md)
   - In-memory storage with maps
   - Save/Get operations
   - Unique ID generation

5. [Package Organization](./step5-refactor-internal-package.md)
   - Split into internal packages
   - Separation of concerns
   - Package structure

6. [State Machine Pattern](./step6-state-machine.md)
   - Status type safety
   - State transition validation
   - Terminal states

7. [Clean Architecture](./step7-clean-architecture.md)
   - Entity/Storage/Use Case layers
   - Dependency injection
   - JobProcessor orchestrator
   - Pure domain entities

8. [Handler Registry](./step8-handler-registry.md)
   - Map job types to handlers
   - RegisterHandler() and GetHandler()
   - Automatic handler lookup
   - Registry pattern

9. [Context Support & Job Retry](./step9-context-and-retry.md)
   - Context.Context for cancellation/timeouts
   - Job retry with MaxRetry and Attempts
   - Exponential backoff
   - Error tracking and timestamps
   - Mark methods (MarkRunning, MarkCompleted, MarkFailed)

10. [Worker Pool & Thread Safety](./step10-worker-pool.md)
   - Worker pool with goroutines
   - Channels for job queue
   - Async job submission
   - Graceful shutdown (context cancellation, WaitGroup)
   - Race condition discovery
   - Thread safety with sync.RWMutex
   - Concurrent map access protection

## Design Principles

- **Start Simple** - Hardcode first, generalize later
- **Test-Driven** - Write tests for each feature
- **Incremental** - Small, focused improvements
- **Separation of Concerns** - Split code logically
- **Type Safety** - Use types over strings
- **Error Handling** - Handle all error cases

## Go Concepts Covered

- Structs and methods
- First-class functions
- Function types
- Error handling and wrapping
- Maps and map lookups
- Package organization
- Internal packages
- Type aliases
- Constants
- State machines
- Testing patterns
- Interfaces
- Dependency injection
- Clean Architecture
- Registry pattern
- Context.Context (cancellation, timeouts)
- time.Time and time.Duration
- Exponential backoff
- Retry logic
- Goroutines (concurrent execution)
- Channels (communication between goroutines)
- Buffered channels
- Select statement
- sync.WaitGroup (waiting for goroutines)
- sync.RWMutex (read-write mutex)
- Worker pool pattern
- Race detector (`go test -race`)
- Graceful shutdown patterns
- Thread safety and concurrent map access
