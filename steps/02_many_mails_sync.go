package steps

import (
	"github.com/MohamedAljoke/goqueue/helpers"
)

func Step02ManyMailsSync() {
	emails := helpers.GenerateRandomEmails(100)

	for _, email := range emails {
		helpers.SendWelcomeEmail(email)
	}

}
