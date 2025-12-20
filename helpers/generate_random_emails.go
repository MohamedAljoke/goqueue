package helpers

import "fmt"

func GenerateRandomEmails(count int) []string {
	emails := []string{}
	for i := 0; i < count; i++ {
		emails = append(emails, fmt.Sprintf("%v-email.mail.com", i))
	}

	return emails
}
