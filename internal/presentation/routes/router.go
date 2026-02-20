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

	// Auth middleware handler function
	auth := middleware.AuthMiddleware(container.JwtService)

	// Create controllers
	authController := controllers.NewAuthController(
		container.RegisterUseCase,
		container.LoginUseCase,
	)

	userController := controllers.NewUserController(
		container.ListUsersUseCase,
		container.GetUserUseCase,
		container.UpdateUserUseCase,
		container.AnonymizeUserUseCase,
	)

	driverController := controllers.NewDriverController(
		container.CreateDriverUseCase,
	)

	brandController := controllers.NewBrandController(
		container.ListBrandsUseCase,
		container.CreateBrandUseCase,
		container.DeleteBrandUseCase,
	)

	colorController := controllers.NewColorController(
		container.ListColorsUseCase,
		container.CreateColorUseCase,
		container.UpdateColorUseCase,
		container.DeleteColorUseCase,
	)

	cityController := controllers.NewCityController(
		container.ListCitiesUseCase,
		container.CreateCityUseCase,
		container.DeleteCityUseCase,
	)

	carController := controllers.NewCarController(
		container.ListCarsUseCase,
		container.CreateCarUseCase,
		container.UpdateCarUseCase,
		container.DeleteCarUseCase,
	)

	tripController := controllers.NewTripController(
		container.ListTripsUseCase,
		container.GetTripUseCase,
		container.FindTripsUseCase,
		container.CreateTripUseCase,
		container.DeleteTripUseCase,
	)

	inscriptionController := controllers.NewInscriptionController(
		container.ListInscriptionsUseCase,
		container.CreateInscriptionUseCase,
		container.DeleteInscriptionUseCase,
		container.ListUserInscriptionsUseCase,
		container.ListTripPassengersUseCase,
	)

	// Register routes under /api/v1
	api := router.Group("/api/v1")

	RegisterAuthRoutes(api, authController)
	RegisterUserRoutes(api, userController, inscriptionController, auth)
	RegisterDriverRoutes(api, driverController, auth)
	RegisterBrandRoutes(api, brandController, auth)
	RegisterColorRoutes(api, colorController, auth)
	RegisterCityRoutes(api, cityController, auth)
	RegisterCarRoutes(api, carController, auth)
	RegisterTripRoutes(api, tripController, inscriptionController, auth)
	RegisterInscriptionRoutes(api, inscriptionController, auth)

	return router
}
