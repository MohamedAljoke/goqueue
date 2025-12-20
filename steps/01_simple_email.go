package steps

import "github.com/MohamedAljoke/goqueue/helpers"

func Step01SimpleEmail() {
	email := "test@example.com"
	helpers.SendWelcomeEmail(email)
}
