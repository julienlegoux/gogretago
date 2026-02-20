package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/presentation/controllers"
	"github.com/lgxju/gogretago/internal/presentation/middleware"
)

// RegisterBrandRoutes registers all brand routes
func RegisterBrandRoutes(router *gin.RouterGroup, brandController *controllers.BrandController, auth gin.HandlerFunc) {
	brands := router.Group("/brands")
	brands.Use(auth)
	brands.GET("", middleware.RequireRole("DRIVER"), brandController.ListBrands)
	brands.POST("", middleware.RequireRole("ADMIN"), brandController.CreateBrand)
	brands.DELETE("/:id", middleware.RequireRole("ADMIN"), brandController.DeleteBrand)
}
