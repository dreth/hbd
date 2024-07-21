package helper

import (
	"errors"
	"regexp"
	"strconv"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if the email has a valid format
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// Check the length of a particular string and return an error if it is too long
func CheckStringLength(field string, str string, maxLength int, minLength int, strictLength int) error {
	if len(str) == 0 {
		return StringInappropriateLengthError(field, "empty")
	} else if len(str) < minLength {
		return StringInappropriateLengthError(field, "too short")
	} else if len(str) > maxLength {
		return StringInappropriateLengthError(field, "too long")
	} else if (strictLength != 0) && (len(str) != strictLength) {
		return StringInappropriateLengthError(field, "incorrect length, should be "+strconv.Itoa(strictLength))
	}
	return nil
}

// Check the length of a particular array of string and return an error if it is too long
func CheckArrayStringLength(field []string, str []string, maxLength []int, minLength []int, strictLength []int) []error {
	errorsList := []error{}

	// Boundary check for array lengths
	if len(field) != len(str) || len(str) != len(maxLength) || len(maxLength) != len(minLength) || len(minLength) != len(strictLength) {
		return append(errorsList, errors.New("array lengths do not match"))
	}

	for i := 0; i < len(str); i++ {
		if len(str[i]) == 0 {
			errorsList = append(errorsList, StringInappropriateLengthError(field[i], "empty"))
		} else if len(str[i]) < minLength[i] {
			errorsList = append(errorsList, StringInappropriateLengthError(field[i], "too short"))
		} else if len(str[i]) > maxLength[i] {
			errorsList = append(errorsList, StringInappropriateLengthError(field[i], "too long"))
		} else if (strictLength[i] != 0) && (len(str[i]) != strictLength[i]) {
			errorsList = append(errorsList, StringInappropriateLengthError(field[i], "incorrect length, should be "+strconv.Itoa(strictLength[i])))
		}
	}
	return errorsList
}
