package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/trip"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// TripController handles trip endpoints
type TripController struct {
	listUseCase   *trip.ListTripsUseCase
	getUseCase    *trip.GetTripUseCase
	findUseCase   *trip.FindTripsUseCase
	createUseCase *trip.CreateTripUseCase
	deleteUseCase *trip.DeleteTripUseCase
}

// NewTripController creates a new TripController
func NewTripController(
	listUseCase *trip.ListTripsUseCase,
	getUseCase *trip.GetTripUseCase,
	findUseCase *trip.FindTripsUseCase,
	createUseCase *trip.CreateTripUseCase,
	deleteUseCase *trip.DeleteTripUseCase,
) *TripController {
	return &TripController{
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
		findUseCase:   findUseCase,
		createUseCase: createUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// ListTrips handles GET /trips
func (ctrl *TripController) ListTrips(c *gin.Context) {
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

// GetTrip handles GET /trips/:id
func (ctrl *TripController) GetTrip(c *gin.Context) {
	id := c.Param("id")

	result, err := ctrl.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// FindTrip handles GET /trips/search
func (ctrl *TripController) FindTrip(c *gin.Context) {
	var query dtos.FindTripQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid query parameters",
			},
		})
		return
	}

	result, err := ctrl.findUseCase.Execute(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CreateTrip handles POST /trips
func (ctrl *TripController) CreateTrip(c *gin.Context) {
	userID := c.GetString("userId")

	var input dtos.CreateTripInput
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeleteTrip handles DELETE /trips/:id
func (ctrl *TripController) DeleteTrip(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	err := ctrl.deleteUseCase.Execute(c.Request.Context(), id, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
