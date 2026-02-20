package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterColorRoutes registers all color routes
func RegisterColorRoutes(router *gin.RouterGroup, colorController *controllers.ColorController, auth gin.HandlerFunc) {
	colors := router.Group("/colors")
	colors.Use(auth)
	colors.GET("", middleware.RequireRole("DRIVER"), colorController.ListColors)
	colors.POST("", middleware.RequireRole("ADMIN"), colorController.CreateColor)
	colors.PATCH("/:id", middleware.RequireRole("ADMIN"), colorController.UpdateColor)
	colors.DELETE("/:id", middleware.RequireRole("ADMIN"), colorController.DeleteColor)
}
