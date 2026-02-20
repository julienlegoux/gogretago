package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterUserRoutes registers all user routes
func RegisterUserRoutes(router *gin.RouterGroup, userController *controllers.UserController, inscriptionController *controllers.InscriptionController, auth gin.HandlerFunc) {
	users := router.Group("/users")
	users.Use(auth)
	{
		users.GET("", middleware.RequireRole("ADMIN"), userController.ListUsers)
		users.GET("/:id", middleware.RequireRole("USER"), userController.GetUser)
		users.PATCH("/me", middleware.RequireRole("USER"), userController.UpdateProfile)
		users.DELETE("/me", middleware.RequireRole("USER"), userController.AnonymizeMe)
		users.DELETE("/:id", middleware.RequireRole("ADMIN"), userController.AnonymizeUser)
		users.GET("/:id/inscriptions", middleware.RequireRole("USER"), inscriptionController.ListUserInscriptions)
	}
}
