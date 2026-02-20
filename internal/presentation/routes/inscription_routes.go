package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterInscriptionRoutes registers all inscription routes
func RegisterInscriptionRoutes(router *gin.RouterGroup, inscriptionController *controllers.InscriptionController, auth gin.HandlerFunc) {
	inscriptions := router.Group("/inscriptions")
	inscriptions.Use(auth)
	{
		inscriptions.GET("", middleware.RequireRole("USER"), inscriptionController.ListInscriptions)
		inscriptions.POST("", middleware.RequireRole("USER"), inscriptionController.CreateInscription)
		inscriptions.DELETE("/:id", middleware.RequireRole("USER"), inscriptionController.DeleteInscription)
	}
}
