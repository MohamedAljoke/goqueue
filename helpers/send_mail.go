package helpers

import (
	"fmt"
	"time"
)

func SendWelcomeEmail(email string) {
	fmt.Println("sending email to:", email)
	time.Sleep(3 * time.Second)
	fmt.Println("email sent to:", email)
}
