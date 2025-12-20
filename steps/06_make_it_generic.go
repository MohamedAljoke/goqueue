package steps

import "fmt"

type HandlerFunc func() error

type Job struct {
	ID     string
	Status string
}

func NewJob() *Job {
	return &Job{
		ID:     "job_12345",
		Status: "pending",
	}
}

func (job *Job) Process(handler HandlerFunc) error {
	job.Status = "processing"
	err := handler()
	if err != nil {
		job.Status = "failed"
		//which processor errored
		return fmt.Errorf("error processing: %w", err)
	}

	job.Status = "completed"

	return nil
}

// func main() {
// 	job := steps.NewJob()

// 	err := job.Process(func() error {
// 		helpers.SendWelcomeEmail("oi")
// 		return nil
// 	})

// 	if err != nil {
// 		panic(err)
// 	}
// }
