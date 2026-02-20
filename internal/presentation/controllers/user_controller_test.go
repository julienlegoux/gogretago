package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/user"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupUserController(t *testing.T) (*UserController, *mocks.MockUserRepository) {
	userRepo := mocks.NewMockUserRepository(t)

	listUC := user.NewListUsersUseCase(userRepo)
	getUC := user.NewGetUserUseCase(userRepo)
	updateUC := user.NewUpdateUserUseCase(userRepo)
	anonymizeUC := user.NewAnonymizeUserUseCase(userRepo)
	ctrl := NewUserController(listUC, getUC, updateUC, anonymizeUC)

	return ctrl, userRepo
}

func TestUserController_ListUsers_Success(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	firstName := "John"
	users := []entities.PublicUser{
		{User: entities.User{ID: "user-1", FirstName: &firstName}, Email: "john@example.com"},
	}
	userRepo.EXPECT().FindAll(mock.Anything).Return(users, nil)

	router := gin.New()
	router.GET("/users", ctrl.ListUsers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestUserController_ListUsers_Error(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	userRepo.EXPECT().FindAll(mock.Anything).Return(nil, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/users", ctrl.ListUsers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_GetUser_OwnProfile(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	firstName := "John"
	userEntity := &entities.PublicUser{
		User:  entities.User{ID: "user-1", FirstName: &firstName},
		Email: "john@example.com",
	}
	userRepo.EXPECT().FindByID(mock.Anything, "user-1").Return(userEntity, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Set("role", "USER")
		c.Next()
	})
	router.GET("/users/:id", ctrl.GetUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/user-1", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestUserController_GetUser_AdminCanViewOther(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	firstName := "Jane"
	userEntity := &entities.PublicUser{
		User:  entities.User{ID: "user-2", FirstName: &firstName},
		Email: "jane@example.com",
	}
	userRepo.EXPECT().FindByID(mock.Anything, "user-2").Return(userEntity, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Set("role", "ADMIN")
		c.Next()
	})
	router.GET("/users/:id", ctrl.GetUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/user-2", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_GetUser_NonAdminCannotViewOther(t *testing.T) {
	ctrl, _ := setupUserController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Set("role", "USER")
		c.Next()
	})
	router.GET("/users/:id", ctrl.GetUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/user-2", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["success"])
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "FORBIDDEN", errObj["code"])
}

func TestUserController_GetUser_DriverCannotViewOther(t *testing.T) {
	ctrl, _ := setupUserController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Set("role", "DRIVER")
		c.Next()
	})
	router.GET("/users/:id", ctrl.GetUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/user-2", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUserController_UpdateProfile_Success(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	firstName := "Jane"
	lastName := "Doe"
	phone := "0612345678"
	existing := &entities.PublicUser{
		User:  entities.User{ID: "user-1"},
		Email: "jane@example.com",
	}
	updated := &entities.PublicUser{
		User:  entities.User{ID: "user-1", FirstName: &firstName, LastName: &lastName, Phone: &phone},
		Email: "jane@example.com",
	}

	userRepo.EXPECT().FindByID(mock.Anything, "user-1").Return(existing, nil)
	userRepo.EXPECT().Update(mock.Anything, "user-1", entities.UpdateUserData{
		FirstName: &firstName,
		LastName:  &lastName,
		Phone:     &phone,
	}).Return(updated, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.PATCH("/users/me", ctrl.UpdateProfile)

	body := `{"firstName":"Jane","lastName":"Doe","phone":"0612345678"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_UpdateProfile_InvalidJSON(t *testing.T) {
	ctrl, _ := setupUserController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.PATCH("/users/me", ctrl.UpdateProfile)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_UpdateProfile_ValidationError(t *testing.T) {
	ctrl, _ := setupUserController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.PATCH("/users/me", ctrl.UpdateProfile)

	body := `{"firstName":"","lastName":"","phone":"123"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_AnonymizeMe_Success(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	existing := &entities.PublicUser{
		User:  entities.User{ID: "user-1"},
		Email: "test@example.com",
	}
	userRepo.EXPECT().FindByID(mock.Anything, "user-1").Return(existing, nil)
	userRepo.EXPECT().Anonymize(mock.Anything, "user-1").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.DELETE("/users/me", ctrl.AnonymizeMe)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/users/me", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUserController_AnonymizeUser_Success(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	existing := &entities.PublicUser{
		User:  entities.User{ID: "user-2"},
		Email: "test@example.com",
	}
	userRepo.EXPECT().FindByID(mock.Anything, "user-2").Return(existing, nil)
	userRepo.EXPECT().Anonymize(mock.Anything, "user-2").Return(nil)

	router := gin.New()
	router.DELETE("/users/:id", ctrl.AnonymizeUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/users/user-2", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUserController_AnonymizeUser_Error(t *testing.T) {
	ctrl, userRepo := setupUserController(t)

	userRepo.EXPECT().FindByID(mock.Anything, "user-2").Return(nil, nil)

	router := gin.New()
	router.DELETE("/users/:id", ctrl.AnonymizeUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/users/user-2", http.NoBody)
	router.ServeHTTP(w, req)

	// c.Error() is called, status not set explicitly
	assert.Equal(t, http.StatusOK, w.Code)
}
