package helper

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if the email has a valid format
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
