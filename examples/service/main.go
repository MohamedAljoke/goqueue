package main

import (
	"context"
	"log"
	"time"

	"github.com/MohamedAljoke/goqueue"
)

// EmailService handles email operations
type EmailService struct {
	smtpHost string
}

func NewEmailService(host string) *EmailService {
	return &EmailService{smtpHost: host}
}

func (s *EmailService) SendWelcomeEmail(ctx context.Context, payload map[string]interface{}) error {
	email := payload["email"]
	name := payload["name"]
	log.Printf("📧 [EmailService] Sending welcome email to %v (%s)", name, email)
	log.Printf("   Using SMTP: %s", s.smtpHost)
	time.Sleep(300 * time.Millisecond)
	return nil
}

func (s *EmailService) SendPasswordReset(ctx context.Context, payload map[string]interface{}) error {
	email := payload["email"]
	log.Printf("🔒 [EmailService] Sending password reset to %v", email)
	time.Sleep(200 * time.Millisecond)
	return nil
}

// PaymentService handles payment operations
type PaymentService struct {
	apiKey string
}

func NewPaymentService(apiKey string) *PaymentService {
	return &PaymentService{apiKey: apiKey}
}

func (s *PaymentService) ChargeCard(ctx context.Context, payload map[string]interface{}) error {
	amount := payload["amount"]
	customerID := payload["customer_id"]
	log.Printf("💳 [PaymentService] Charging $%v to customer %v", amount, customerID)
	log.Printf("   Using API key: %s...", s.apiKey[:10])
	time.Sleep(500 * time.Millisecond)
	return nil
}

// AnalyticsService handles analytics
type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

func (s *AnalyticsService) TrackEvent(ctx context.Context, payload map[string]interface{}) error {
	event := payload["event"]
	userID := payload["user_id"]
	log.Printf("📊 [AnalyticsService] Tracking event '%v' for user %v", event, userID)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func main() {
	log.Println("=== GoQueue Service Example ===")
	log.Println("Demonstrating how to integrate GoQueue with your existing services")

	// Initialize your services
	emailService := NewEmailService("smtp.example.com")
	paymentService := NewPaymentService("sk_test_abc123xyz")
	analyticsService := NewAnalyticsService()

	// Create queue
	q := goqueue.New(
		goqueue.WithWorkers(4),
		goqueue.WithBufferSize(15),
	)

	// Register handlers - connecting job types to service methods
	q.RegisterHandler("email.welcome", emailService.SendWelcomeEmail)
	q.RegisterHandler("email.password_reset", emailService.SendPasswordReset)
	q.RegisterHandler("payment.charge", paymentService.ChargeCard)
	q.RegisterHandler("analytics.track", analyticsService.TrackEvent)

	// Start queue
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go q.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	log.Println("\n--- Submitting Jobs ---")

	// User registration flow
	q.Submit(ctx, "email.welcome", map[string]interface{}{
		"email": "alice@example.com",
		"name":  "Alice",
	}, 3)

	q.Submit(ctx, "analytics.track", map[string]interface{}{
		"event":   "user.registered",
		"user_id": "user_123",
	}, 2)

	// Payment flow
	q.Submit(ctx, "payment.charge", map[string]interface{}{
		"amount":      99.99,
		"customer_id": "cus_abc123",
	}, 5)

	q.Submit(ctx, "analytics.track", map[string]interface{}{
		"event":   "payment.completed",
		"user_id": "user_123",
	}, 2)

	// Password reset flow
	q.Submit(ctx, "email.password_reset", map[string]interface{}{
		"email": "bob@example.com",
	}, 3)

	// Let jobs process
	time.Sleep(3 * time.Second)

	log.Println("\n--- Shutting Down ---")
	cancel()
	time.Sleep(500 * time.Millisecond)
	log.Println("✅ All services stopped gracefully")
}
