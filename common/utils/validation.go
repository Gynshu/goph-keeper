package utils

import "net/mail"

// ValidateEmail checks if an email address is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
