package main

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"hbd/auth"
	"hbd/env"
	"hbd/middlewares"
)

func main() {

	boil.SetDB(env.DB)

	router := gin.Default()

	// Apply middleware
	router.Use(middlewares.RateLimitMiddleware())

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.DELETE("/delete-user", auth.DeleteUser)
	router.PATCH("/modify-user", auth.ModifyUser)
	router.GET("/generate-encryption-key", auth.GetEncryptionKey)

	router.Run(":8080")
}
