package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
)

// RegisterAuthRoutes registers all auth routes
func RegisterAuthRoutes(router *gin.RouterGroup, authController *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}
}
