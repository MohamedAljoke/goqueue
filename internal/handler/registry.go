package handler

import (
	"fmt"
	"sync"

	"github.com/MohamedAljoke/goqueue/internal/entity"
)

type Registry struct {
	handlers map[string]entity.HandlerFunc
	mu       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]entity.HandlerFunc),
	}
}

func (r *Registry) RegisterHandler(jobType string, handler entity.HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[jobType] = handler

}

func (r *Registry) GetHandler(jobType string) (entity.HandlerFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, exists := r.handlers[jobType]
	if !exists {
		return nil, fmt.Errorf("no handler registered for job type: %s", jobType)
	}
	return handler, nil
}
