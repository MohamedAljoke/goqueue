package main

import (
	"github.com/MohamedAljoke/goqueue/helpers"
	"github.com/MohamedAljoke/goqueue/steps"
)

func main() {
	job := steps.NewJob()

	err := job.Process(func() error {
		helpers.SendWelcomeEmail("oi")
		return nil
	})

	if err != nil {
		panic(err)
	}
}
