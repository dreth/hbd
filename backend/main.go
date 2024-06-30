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

	docs "hbd/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.Use(middlewares.SwaggerHostMiddleware())

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "HBD API"
	docs.SwaggerInfo.Description = "HBD endpoints for the HBD application frontend"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.DELETE("/delete-user", auth.DeleteUser)
	router.PUT("/modify-user", auth.ModifyUser)
	router.GET("/generate-encryption-key", auth.GetEncryptionKey)

	// Birthday routes
	router.POST("/check-birthdays", birthdays.CallReminderChecker)
	router.POST("/birthday", birthdays.AddBirthday)
	router.PUT("/birthday", birthdays.ModifyBirthday)
	router.DELETE("/birthday", birthdays.DeleteBirthday)

	// Run the server
	router.Run(":8417")
}
