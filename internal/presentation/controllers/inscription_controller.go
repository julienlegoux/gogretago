package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/inscription"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// InscriptionController handles inscription endpoints
type InscriptionController struct {
	listUseCase                *inscription.ListInscriptionsUseCase
	createUseCase              *inscription.CreateInscriptionUseCase
	deleteUseCase              *inscription.DeleteInscriptionUseCase
	listUserInscriptionsUseCase *inscription.ListUserInscriptionsUseCase
	listTripPassengersUseCase   *inscription.ListTripPassengersUseCase
}

// NewInscriptionController creates a new InscriptionController
func NewInscriptionController(
	listUseCase *inscription.ListInscriptionsUseCase,
	createUseCase *inscription.CreateInscriptionUseCase,
	deleteUseCase *inscription.DeleteInscriptionUseCase,
	listUserInscriptionsUseCase *inscription.ListUserInscriptionsUseCase,
	listTripPassengersUseCase *inscription.ListTripPassengersUseCase,
) *InscriptionController {
	return &InscriptionController{
		listUseCase:                listUseCase,
		createUseCase:              createUseCase,
		deleteUseCase:              deleteUseCase,
		listUserInscriptionsUseCase: listUserInscriptionsUseCase,
		listTripPassengersUseCase:   listTripPassengersUseCase,
	}
}

// ListInscriptions handles GET /inscriptions
func (ctrl *InscriptionController) ListInscriptions(c *gin.Context) {
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

// CreateInscription handles POST /inscriptions
func (ctrl *InscriptionController) CreateInscription(c *gin.Context) {
	userID := c.GetString("userId")

	var input dtos.CreateInscriptionInput
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

// DeleteInscription handles DELETE /inscriptions/:id
func (ctrl *InscriptionController) DeleteInscription(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	err := ctrl.deleteUseCase.Execute(c.Request.Context(), id, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUserInscriptions handles GET /users/:id/inscriptions
func (ctrl *InscriptionController) ListUserInscriptions(c *gin.Context) {
	userID := c.Param("id")

	result, err := ctrl.listUserInscriptionsUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListTripPassengers handles GET /trips/:id/passengers
func (ctrl *InscriptionController) ListTripPassengers(c *gin.Context) {
	tripID := c.Param("id")

	result, err := ctrl.listTripPassengersUseCase.Execute(c.Request.Context(), tripID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
