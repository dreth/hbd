package helper

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Helper function to handle errors and send JSON response
func HE(c *gin.Context, err error, statusCode int, message string, useDefaultErrorMessage bool) bool {
	if err != nil {
		if !useDefaultErrorMessage {
			log.Println("Error:", err)
			c.JSON(statusCode, gin.H{"error": message})
		} else {
			log.Println("Error:", err.Error())
			c.JSON(statusCode, gin.H{"error": "An error occurred"})
		}
		return true
	}
	return false
}
