package helper

import (
	"hbd/structs"
	"log"

	"github.com/gin-gonic/gin"
)

// Helper function to handle errors and send JSON response
func HE(c *gin.Context, err error, statusCode int, message string, useDefaultErrorMessage bool) bool {
	if err != nil {
		var errorMessage string
		if !useDefaultErrorMessage {
			errorMessage = message
		} else {
			errorMessage = "An error occurred"
		}
		log.Println("Error:", err)
		c.JSON(statusCode, structs.Error{Error: errorMessage})
		return true
	}
	return false
}
