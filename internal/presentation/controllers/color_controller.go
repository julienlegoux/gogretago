package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/color"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// ColorController handles color endpoints
type ColorController struct {
	listUseCase   *color.ListColorsUseCase
	createUseCase *color.CreateColorUseCase
	updateUseCase *color.UpdateColorUseCase
	deleteUseCase *color.DeleteColorUseCase
}

// NewColorController creates a new ColorController
func NewColorController(
	listUseCase *color.ListColorsUseCase,
	createUseCase *color.CreateColorUseCase,
	updateUseCase *color.UpdateColorUseCase,
	deleteUseCase *color.DeleteColorUseCase,
) *ColorController {
	return &ColorController{
		listUseCase:   listUseCase,
		createUseCase: createUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// ListColors handles GET /colors
func (ctrl *ColorController) ListColors(c *gin.Context) {
	params := parsePagination(c)

	result, err := ctrl.listUseCase.Execute(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result.Data,
		"meta":    result.Meta,
	})
}

// CreateColor handles POST /colors
func (ctrl *ColorController) CreateColor(c *gin.Context) {
	var input dtos.CreateColorInput
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateColor handles PATCH /colors/:id
func (ctrl *ColorController) UpdateColor(c *gin.Context) {
	id := c.Param("id")

	var input dtos.UpdateColorInput
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
	result, err := ctrl.updateUseCase.Execute(c.Request.Context(), id, input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeleteColor handles DELETE /colors/:id
func (ctrl *ColorController) DeleteColor(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
