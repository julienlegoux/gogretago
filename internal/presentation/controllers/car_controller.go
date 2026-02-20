package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/car"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// CarController handles car endpoints
type CarController struct {
	listUseCase   *car.ListCarsUseCase
	createUseCase *car.CreateCarUseCase
	updateUseCase *car.UpdateCarUseCase
	deleteUseCase *car.DeleteCarUseCase
}

// NewCarController creates a new CarController
func NewCarController(
	listUseCase *car.ListCarsUseCase,
	createUseCase *car.CreateCarUseCase,
	updateUseCase *car.UpdateCarUseCase,
	deleteUseCase *car.DeleteCarUseCase,
) *CarController {
	return &CarController{
		listUseCase:   listUseCase,
		createUseCase: createUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// ListCars handles GET /cars
func (ctrl *CarController) ListCars(c *gin.Context) {
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

// CreateCar handles POST /cars
func (ctrl *CarController) CreateCar(c *gin.Context) {
	userID := c.GetString("userId")

	var input dtos.CreateCarInput
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
	result, err := ctrl.createUseCase.Execute(c.Request.Context(), userID, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateCar handles PUT /cars/:id
func (ctrl *CarController) UpdateCar(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	var input dtos.UpdateCarInput
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

	// Convert to UpdateCarData (all fields set for PUT)
	data := car.UpdateCarData{
		Model:        &input.Model,
		BrandID:      &input.BrandID,
		LicensePlate: &input.LicensePlate,
	}

	// Execute use case
	result, err := ctrl.updateUseCase.Execute(c.Request.Context(), id, userID, data)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// PatchCar handles PATCH /cars/:id
func (ctrl *CarController) PatchCar(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	var input dtos.PatchCarInput
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

	// Convert to UpdateCarData (only set fields for PATCH)
	data := car.UpdateCarData{
		Model:        input.Model,
		BrandID:      input.BrandID,
		LicensePlate: input.LicensePlate,
	}

	// Execute use case
	result, err := ctrl.updateUseCase.Execute(c.Request.Context(), id, userID, data)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeleteCar handles DELETE /cars/:id
func (ctrl *CarController) DeleteCar(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	err := ctrl.deleteUseCase.Execute(c.Request.Context(), id, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
