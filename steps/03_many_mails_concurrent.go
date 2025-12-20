package steps

import (
	"sync"

	"github.com/MohamedAljoke/goqueue/helpers"
)

func Step03ManyMails() {
	emails := helpers.GenerateRandomEmails(100)

	wg := sync.WaitGroup{}
	for _, email := range emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			helpers.SendWelcomeEmail(email)
		}(email)
	}

	wg.Wait()
}
