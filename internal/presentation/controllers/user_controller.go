package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/user"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// UserController handles user endpoints
type UserController struct {
	listUseCase      *user.ListUsersUseCase
	getUseCase       *user.GetUserUseCase
	updateUseCase    *user.UpdateUserUseCase
	anonymizeUseCase *user.AnonymizeUserUseCase
}

// NewUserController creates a new UserController
func NewUserController(
	listUseCase *user.ListUsersUseCase,
	getUseCase *user.GetUserUseCase,
	updateUseCase *user.UpdateUserUseCase,
	anonymizeUseCase *user.AnonymizeUserUseCase,
) *UserController {
	return &UserController{
		listUseCase:      listUseCase,
		getUseCase:       getUseCase,
		updateUseCase:    updateUseCase,
		anonymizeUseCase: anonymizeUseCase,
	}
}

// ListUsers handles GET /users
func (ctrl *UserController) ListUsers(c *gin.Context) {
	result, err := ctrl.listUseCase.Execute(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetUser handles GET /users/:id
// Users can only view their own profile unless they have ADMIN role.
func (ctrl *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	requestingUserID := c.GetString("userId")
	role := c.GetString("role")

	if id != requestingUserID && role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "Insufficient permissions",
			},
		})
		return
	}

	result, err := ctrl.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateProfile handles PATCH /users/me
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userId")

	var input dtos.UpdateProfileInput
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
	result, err := ctrl.updateUseCase.Execute(c.Request.Context(), userID, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AnonymizeMe handles DELETE /users/me
func (ctrl *UserController) AnonymizeMe(c *gin.Context) {
	userID := c.GetString("userId")

	err := ctrl.anonymizeUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// AnonymizeUser handles DELETE /users/:id
func (ctrl *UserController) AnonymizeUser(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.anonymizeUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
