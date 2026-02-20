package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterCityRoutes registers all city routes
func RegisterCityRoutes(router *gin.RouterGroup, cityController *controllers.CityController, auth gin.HandlerFunc) {
	cities := router.Group("/cities")
	cities.Use(auth)
	{
		cities.GET("", middleware.RequireRole("USER"), cityController.ListCities)
		cities.POST("", middleware.RequireRole("USER"), cityController.CreateCity)
		cities.DELETE("/:id", middleware.RequireRole("ADMIN"), cityController.DeleteCity)
	}
}
