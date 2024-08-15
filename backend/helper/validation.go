package helper

import (
	"errors"
	"regexp"
	"strconv"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9]{2,}$`)

// IsValidEmail checks if the email has a valid format
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// Check the length of a particular string and return an error if it is too long
func CheckStringLength(field string, str string, maxLength int, minLength int, strictLength int, allowEmpty bool) error {
	if len(str) == 0 {
		if !allowEmpty {
			return StringInappropriateLengthError(field, "empty")
		}
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
func CheckArrayStringLength(field []string, str []string, maxLength []int, minLength []int, strictLength []int, allowEmpty []bool) []error {
	errorsList := []error{}

	// Boundary check for array lengths
	if len(field) != len(str) || len(str) != len(maxLength) || len(maxLength) != len(minLength) || len(minLength) != len(strictLength) {
		return append(errorsList, errors.New("array lengths do not match"))
	}

	for i := 0; i < len(str); i++ {
		err := CheckStringLength(field[i], str[i], maxLength[i], minLength[i], strictLength[i], allowEmpty[i])
		errorsList = append(errorsList, err)
	}
	return errorsList
}
