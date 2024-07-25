package main

import (
	"time"

	"github.com/gin-contrib/cors"
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
	c.AddFunc("* * * * *", birthdays.CheckReminders)
	c.Start()

	// Initialize the database connection and run migrations
	boil.SetDB(env.DB)
	db.RunMigrations(env.DB)

	// Initialize the Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8417", "http://localhost:8418", "http://localhost:3000", "http://0.0.0.0:8418", env.CDFE, env.CDBE},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

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
	router.GET("/generate-password", auth.GetPassword)

	// Requires authentication
	authenticated := router.Group("/")
	authenticated.Use(middlewares.JWTAuthMiddleware())
	{
		// User routes
		authenticated.GET("/me", auth.Me)
		authenticated.DELETE("/delete-user", auth.DeleteUser)
		authenticated.PUT("/modify-user", auth.ModifyUser)

		// Birthday routes
		authenticated.PATCH("/check-birthdays", birthdays.CallReminderChecker)
		authenticated.POST("/add-birthday", birthdays.AddBirthday)
		authenticated.PUT("/modify-birthday", birthdays.ModifyBirthday)
		authenticated.DELETE("/delete-birthday", birthdays.DeleteBirthday)
	}

	// Run the server
	router.Run("0.0.0.0:8417")
}
