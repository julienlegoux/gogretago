package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/city"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// CityController handles city endpoints
type CityController struct {
	listUseCase   *city.ListCitiesUseCase
	createUseCase *city.CreateCityUseCase
	deleteUseCase *city.DeleteCityUseCase
}

// NewCityController creates a new CityController
func NewCityController(
	listUseCase *city.ListCitiesUseCase,
	createUseCase *city.CreateCityUseCase,
	deleteUseCase *city.DeleteCityUseCase,
) *CityController {
	return &CityController{
		listUseCase:   listUseCase,
		createUseCase: createUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// ListCities handles GET /cities
func (ctrl *CityController) ListCities(c *gin.Context) {
	params := parsePagination(c)

	result, err := ctrl.listUseCase.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result.Data,
		"meta":    result.Meta,
	})
}

// CreateCity handles POST /cities
func (ctrl *CityController) CreateCity(c *gin.Context) {
	var input dtos.CreateCityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request body",
			},
		})
		return
	}

	// Validate input
	validate := validators.GetValidator()
	if err := validate.Struct(input); err != nil {
		details := validators.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Validation failed",
				"details": details,
			},
		})
		return
	}

	// Execute use case
	result, err := ctrl.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeleteCity handles DELETE /cities/:id
func (ctrl *CityController) DeleteCity(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
