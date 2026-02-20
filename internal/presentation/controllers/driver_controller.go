package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/driver"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// DriverController handles driver endpoints
type DriverController struct {
	createUseCase *driver.CreateDriverUseCase
}

// NewDriverController creates a new DriverController
func NewDriverController(
	createUseCase *driver.CreateDriverUseCase,
) *DriverController {
	return &DriverController{
		createUseCase: createUseCase,
	}
}

// CreateDriver handles POST /drivers
func (ctrl *DriverController) CreateDriver(c *gin.Context) {
	userID := c.GetString("userId")

	var input dtos.CreateDriverInput
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
