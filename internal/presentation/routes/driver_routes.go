package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterDriverRoutes registers all driver routes
func RegisterDriverRoutes(router *gin.RouterGroup, driverController *controllers.DriverController, auth gin.HandlerFunc) {
	drivers := router.Group("/drivers")
	drivers.Use(auth)
	drivers.POST("", middleware.RequireRole("USER"), driverController.CreateDriver)
}
