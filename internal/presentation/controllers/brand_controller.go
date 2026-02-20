package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/brand"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// BrandController handles brand endpoints
type BrandController struct {
	listUseCase   *brand.ListBrandsUseCase
	createUseCase *brand.CreateBrandUseCase
	deleteUseCase *brand.DeleteBrandUseCase
}

// NewBrandController creates a new BrandController
func NewBrandController(
	listUseCase *brand.ListBrandsUseCase,
	createUseCase *brand.CreateBrandUseCase,
	deleteUseCase *brand.DeleteBrandUseCase,
) *BrandController {
	return &BrandController{
		listUseCase:   listUseCase,
		createUseCase: createUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// ListBrands handles GET /brands
func (ctrl *BrandController) ListBrands(c *gin.Context) {
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

// CreateBrand handles POST /brands
func (ctrl *BrandController) CreateBrand(c *gin.Context) {
	var input dtos.CreateBrandInput
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

// DeleteBrand handles DELETE /brands/:id
func (ctrl *BrandController) DeleteBrand(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
