package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/application/usecases/auth"
	"github.com/lgxju/gogretago/internal/presentation/validators"
)

// AuthController handles authentication endpoints
type AuthController struct {
	registerUseCase *auth.RegisterUseCase
	loginUseCase    *auth.LoginUseCase
}

// NewAuthController creates a new AuthController
func NewAuthController(
	registerUseCase *auth.RegisterUseCase,
	loginUseCase *auth.LoginUseCase,
) *AuthController {
	return &AuthController{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
	}
}

// Register handles user registration
func (ctrl *AuthController) Register(c *gin.Context) {
	var input dtos.RegisterInput

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
	result, err := ctrl.registerUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

// Login handles user login
func (ctrl *AuthController) Login(c *gin.Context) {
	var input dtos.LoginInput

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
	result, err := ctrl.loginUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
