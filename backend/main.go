package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"hbd/auth"
	"hbd/birthdays"
	"hbd/db"
	"hbd/env"
	"hbd/middlewares"
)

func main() {
	// Set up the cron job to check for birthday reminder checks every minute
	c := cron.New()
	c.AddFunc("*/1 * * * *", birthdays.CheckReminders)
	c.Start()

	// Initialize the database connection and run migrations
	boil.SetDB(env.DB)
	db.RunMigrations(env.DB)

	// Initialize the Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middlewares.RateLimitMiddleware())

	// Routes
	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.DELETE("/delete-user", auth.DeleteUser)
	router.PUT("/modify-user", auth.ModifyUser)
	router.GET("/generate-encryption-key", auth.GetEncryptionKey)
	router.GET("/birthdays", birthdays.CallReminderChecker)

	// Run the server
	router.Run(":8417")
}
