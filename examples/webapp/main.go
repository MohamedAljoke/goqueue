package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MohamedAljoke/goqueue"
)

// App holds application dependencies
type App struct {
	queue *goqueue.Queue
}

func main() {
	log.Println("=== GoQueue Web App Example ===")

	// Create queue with custom configuration
	q := goqueue.New(
		goqueue.WithWorkers(5),
		goqueue.WithBufferSize(20),
	)

	// Register business logic handlers
	q.RegisterHandler("send_notification", sendNotificationHandler)
	q.RegisterHandler("process_order", processOrderHandler)
	q.RegisterHandler("generate_invoice", generateInvoiceHandler)

	app := &App{queue: q}

	// Start queue workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go q.Start(ctx)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", app.submitJobHandler)
	mux.HandleFunc("/jobs/status", app.jobStatusHandler)

	// Start HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("🌐 Server listening on :8080")
		log.Println("Try: curl -X POST http://localhost:8080/jobs -d '{\"type\":\"send_notification\",\"payload\":{\"user\":\"alice\"}}'")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	cancel()
	server.Shutdown(context.Background())
	time.Sleep(1 * time.Second)
}

// submitJobHandler handles POST /jobs
func (app *App) submitJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type     string                 `json:"type"`
		Payload  map[string]interface{} `json:"payload"`
		MaxRetry int                    `json:"max_retry"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.MaxRetry == 0 {
		req.MaxRetry = 3
	}

	jobID, err := app.queue.Submit(r.Context(), req.Type, req.Payload, req.MaxRetry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status": "queued",
	})
}

// jobStatusHandler handles GET /jobs/status?id=xxx
func (app *App) jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	job, err := app.queue.GetJob(r.Context(), jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// Business logic handlers
func sendNotificationHandler(ctx context.Context, payload map[string]interface{}) error {
	user := payload["user"]
	log.Printf("📧 Sending notification to user: %v", user)
	time.Sleep(500 * time.Millisecond)
	return nil
}

func processOrderHandler(ctx context.Context, payload map[string]interface{}) error {
	orderID := payload["order_id"]
	log.Printf("📦 Processing order: %v", orderID)
	time.Sleep(1 * time.Second)
	return nil
}

func generateInvoiceHandler(ctx context.Context, payload map[string]interface{}) error {
	customerID := payload["customer_id"]
	log.Printf("🧾 Generating invoice for customer: %v", customerID)
	time.Sleep(800 * time.Millisecond)
	return nil
}
