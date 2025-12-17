package examples

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// SendEmailHandler simulates sending emails
func SendEmailHandler(ctx context.Context, payload map[string]interface{}) error {
	to := payload["to"]
	subject := payload["subject"]

	// Simulate work
	time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)

	// Simulate random failures (20% chance)
	if rand.Float32() < 0.2 {
		return fmt.Errorf("SMTP connection failed")
	}

	log.Printf("  ✉️  Email sent to %v: %v", to, subject)
	return nil
}

// ProcessFileHandler simulates file processing
func ProcessFileHandler(ctx context.Context, payload map[string]interface{}) error {
	filename := payload["filename"]

	// Simulate work
	time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)

	// Simulate random failures (10% chance)
	if rand.Float32() < 0.1 {
		return fmt.Errorf("file corrupted: %v", filename)
	}

	log.Printf("  📁 File processed: %v", filename)
	return nil
}

// GenerateReportHandler simulates report generation
func GenerateReportHandler(ctx context.Context, payload map[string]interface{}) error {
	reportType := payload["type"]

	// Simulate work
	time.Sleep(time.Duration(rand.Intn(800)+300) * time.Millisecond)

	log.Printf("  📊 Report generated: %v", reportType)
	return nil
}
