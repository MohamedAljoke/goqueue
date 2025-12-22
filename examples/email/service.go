package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MohamedAljoke/goqueue"
)

func main() {
	// Create queue
	q := goqueue.NewQueue()

	// Register handler
	q.RegisterHandler("send_welcome_email", func(
		ctx context.Context,
		payload map[string]any,
	) error {
		emailAddr := payload["email"].(string)
		return SendWelcome(emailAddr)
	})

	// Start workers
	q.Start()

	// Simulate app events
	submitUsers(q)

	fmt.Print("teste1")
	// Graceful shutdown (SIGTERM / CTRL+C)
	waitForShutdown(q)
	fmt.Print("teste2")
}
func SendWelcome(email string) error {
	fmt.Println("sending welcome email to:", email)

	// simulate external dependency
	time.Sleep(500 * time.Millisecond)

	// simulate occasional failure
	if email == "fail@test.com" {
		return errors.New("smtp error")
	}

	fmt.Println("email sent to:", email)
	return nil
}

func submitUsers(q *goqueue.Queue) {
	users := []string{
		"user1@test.com",
		"fail@test.com",
		"user2@test.com",
	}

	for _, email := range users {
		job, err := q.SubmitJob(
			context.Background(),
			"send_welcome_email",
			map[string]any{
				"email": email,
			},
			3, // max retries
		)

		if err != nil {
			fmt.Println("failed to submit job:", err)
			continue
		}

		fmt.Println("job submitted:", job.ID)
	}
}
func waitForShutdown(q *goqueue.Queue) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	fmt.Println("shutting down queue...")
	q.Stop()
	fmt.Println("queue stopped")
}
