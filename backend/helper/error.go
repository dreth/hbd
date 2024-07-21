package helper

import (
	"errors"
	"hbd/structs"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// Helper function to handle errors and send JSON response
func HE(c *gin.Context, err error, statusCode int, message string, useDefaultErrorMessage bool) bool {
	if err != nil {
		var errorMessage string
		if !useDefaultErrorMessage {
			errorMessage = message
		} else {
			errorMessage = err.Error()
		}
		log.Println("Error:", err)
		c.JSON(statusCode, structs.Error{Error: errorMessage})
		return true
	}
	return false
}

// Loop over elements in the error array, if they're all nil, return nil
// Otherwise, return the error array with the non-nil errors
func CheckErrors(errorsList []error) []error {
	for _, err := range errorsList {
		if err != nil {
			return errorsList
		}
	}
	return nil
}

// Concatenate the string of all the errors in the error array if they're not nil
func ConcatenateErrors(errorsList []error) string {
	var errorsArray []string
	for _, err := range errorsList {
		if err != nil {
			errorsArray = append(errorsArray, err.Error())
		}
	}
	return strings.Join(errorsArray, ", ")
}

// Possible errors
func StringInappropriateLengthError(field string, message string) error {
	return errors.New("Field '" + field + "' is " + message)
}
