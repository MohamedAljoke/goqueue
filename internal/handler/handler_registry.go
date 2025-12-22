package handler

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type HandlerFunc func(ctx context.Context, payload map[string]any) error

type HandlerRegistry struct {
	handlers map[string]HandlerFunc
	mu       sync.RWMutex
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *HandlerRegistry) Register(jobType string, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[jobType] = handler
}

var ErrRegistryNotFound = errors.New("handler registered for job type")

func (r *HandlerRegistry) Get(jobType string) (HandlerFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[jobType]

	if !ok {
		return nil, fmt.Errorf(
			"%w: %s ",
			ErrRegistryNotFound,
			jobType,
		)
	}

	return h, nil
}
