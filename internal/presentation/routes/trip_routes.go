package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterTripRoutes registers all trip routes
func RegisterTripRoutes(router *gin.RouterGroup, tripController *controllers.TripController, inscriptionController *controllers.InscriptionController, auth gin.HandlerFunc) {
	trips := router.Group("/trips")
	trips.Use(auth)
	trips.GET("", middleware.RequireRole("USER"), tripController.ListTrips)
	trips.GET("/search", middleware.RequireRole("USER"), tripController.FindTrip)
	trips.GET("/:id", middleware.RequireRole("USER"), tripController.GetTrip)
	trips.POST("", middleware.RequireRole("DRIVER"), tripController.CreateTrip)
	trips.DELETE("/:id", middleware.RequireRole("DRIVER"), tripController.DeleteTrip)
	trips.GET("/:id/passengers", middleware.RequireRole("USER"), inscriptionController.ListTripPassengers)
}
