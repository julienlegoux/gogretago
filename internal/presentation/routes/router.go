package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/infrastructure/di"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(container *di.Container) *gin.Engine {
	router := gin.Default()

	// Apply global error handler
	router.Use(middleware.ErrorHandler())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Create controllers
	authController := controllers.NewAuthController(
		container.RegisterUseCase,
		container.LoginUseCase,
	)

	// Register routes
	api := router.Group("")
	RegisterAuthRoutes(api, authController)

	return router
}
