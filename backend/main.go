package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"hbd/auth"
	"hbd/birthdays"
	"hbd/env"
	"hbd/middlewares"
)

func main() {
	c := cron.New()
	c.AddFunc("*/1 * * * *", birthdays.CheckReminders)
	c.Start()

	boil.SetDB(env.DB)

	router := gin.Default()

	// Apply middleware
	router.Use(middlewares.RateLimitMiddleware())

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.DELETE("/delete-user", auth.DeleteUser)
	router.PUT("/modify-user", auth.ModifyUser)
	router.GET("/generate-encryption-key", auth.GetEncryptionKey)
	router.GET("/birthdays", birthdays.CallReminderChecker)

	router.Run(":8080")
}
