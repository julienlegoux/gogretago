package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterAuthRoutes registers all auth routes with rate limiting
func RegisterAuthRoutes(router *gin.RouterGroup, authController *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimiter(3), authController.Register) // 3 req/min
		auth.POST("/login", middleware.RateLimiter(5), authController.Login)       // 5 req/min
	}
}
