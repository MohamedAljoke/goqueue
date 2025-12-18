package handler_test

import (
	"context"
	"sync"
	"testing"

	"github.com/MohamedAljoke/goqueue/internal/handler"
)

// This test demonstrates the race condition when multiple goroutines
// access the registry concurrently WITHOUT mutex protection
func TestRegistryConcurrentAccess(t *testing.T) {
	registry := handler.NewRegistry()

	// Register initial handler
	registry.RegisterHandler("email", func(ctx context.Context, payload map[string]interface{}) error {
		return nil
	})

	var wg sync.WaitGroup

	// Simulate 10 workers reading handlers concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				registry.GetHandler("email")
			}
		}()
	}

	// Simulate 5 goroutines registering new handlers concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				jobType := "sms"
				registry.RegisterHandler(jobType, func(ctx context.Context, payload map[string]interface{}) error {
					return nil
				})
			}
		}(i)
	}

	wg.Wait()
}
