package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MohamedAljoke/goqueue"
)

func main() {
	log.Println("=== Simple GoQueue Example ===")

	// Create a queue with default options
	q := goqueue.New()

	// Register your handlers
	q.RegisterHandler("greet", greetHandler)
	q.RegisterHandler("calculate", calculateHandler)

	// Start processing in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go q.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Submit some jobs
	jobID1, _ := q.Submit(ctx, "greet", map[string]interface{}{
		"name": "Alice",
	}, 3)

	jobID2, _ := q.Submit(ctx, "calculate", map[string]interface{}{
		"operation": "add",
		"a":         10,
		"b":         5,
	}, 3)

	log.Printf("Submitted jobs: %s, %s", jobID1, jobID2)

	// Wait for processing
	time.Sleep(2 * time.Second)

	// Check job status
	job, _ := q.GetJob(ctx, jobID1)
	log.Printf("Job %s status: %s", job.ID, job.Status)

	// Shutdown
	cancel()
	time.Sleep(500 * time.Millisecond)
}

// greetHandler is your custom business logic
func greetHandler(ctx context.Context, payload map[string]interface{}) error {
	name := payload["name"]
	log.Printf("👋 Hello, %s!", name)
	return nil
}

// calculateHandler demonstrates another handler
func calculateHandler(ctx context.Context, payload map[string]interface{}) error {
	op := payload["operation"]
	a := payload["a"].(int)
	b := payload["b"].(int)

	var result int
	switch op {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	default:
		return fmt.Errorf("unknown operation: %v", op)
	}

	log.Printf("🧮 %d %s %d = %d", a, op, b, result)
	return nil
}
