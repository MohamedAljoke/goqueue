package handler

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
)

// Func is a function that processes a job payload
type Func func(ctx context.Context, payload map[string]interface{}) error

// Registry manages job type handlers
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]Func
}

// NewRegistry creates a new handler registry
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Func),
	}
}

// Register adds a handler for a job type
func (r *Registry) Register(jobType string, fn Func) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[jobType] = fn
}

// Get retrieves a handler for a job type
func (r *Registry) Get(jobType string) (Func, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, ok := r.handlers[jobType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrHandlerNotFound, jobType)
	}

	return fn, nil
}

// Has checks if a handler exists for a job type
func (r *Registry) Has(jobType string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.handlers[jobType]
	return ok
}
