package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterCarRoutes registers all car routes
func RegisterCarRoutes(router *gin.RouterGroup, carController *controllers.CarController, auth gin.HandlerFunc) {
	cars := router.Group("/cars")
	cars.Use(auth)
	cars.GET("", middleware.RequireRole("DRIVER"), carController.ListCars)
	cars.POST("", middleware.RequireRole("DRIVER"), carController.CreateCar)
	cars.PUT("/:id", middleware.RequireRole("DRIVER"), carController.UpdateCar)
	cars.PATCH("/:id", middleware.RequireRole("DRIVER"), carController.PatchCar)
	cars.DELETE("/:id", middleware.RequireRole("DRIVER"), carController.DeleteCar)
}
