package helper

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Helper function to handle errors and send JSON response
func HE(c *gin.Context, err error, statusCode int, message string) bool {
	if err != nil {
		log.Println("Error:", err)
		c.JSON(statusCode, gin.H{"error": message})
		return true
	}
	return false
}
