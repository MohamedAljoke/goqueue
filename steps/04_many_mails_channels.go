package steps

import (
	"fmt"
	"sync"
	"time"

	"github.com/MohamedAljoke/goqueue/helpers"
)

func Step04ManyMailsChannels(mails int) {
	workers := mails / 2
	emailChan := make(chan string)

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for email := range emailChan {
				helpers.SendWelcomeEmail(email)
			}
		}()
	}

	go func() {
		emails := helpers.GenerateRandomEmails(mails)
		for _, email := range emails {
			emailChan <- email
			time.Sleep(100 * time.Millisecond)
		}
		close(emailChan)
	}()

	wg.Wait()
	fmt.Println("teste")

}
