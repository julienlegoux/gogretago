package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/driver"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupDriverController(t *testing.T) (
	*DriverController,
	*mocks.MockDriverRepository,
	*mocks.MockUserRepository,
	*mocks.MockAuthRepository,
) {
	driverRepo := mocks.NewMockDriverRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	authRepo := mocks.NewMockAuthRepository(t)

	createUC := driver.NewCreateDriverUseCase(driverRepo, userRepo, authRepo)
	ctrl := NewDriverController(createUC)

	return ctrl, driverRepo, userRepo, authRepo
}

func TestDriverController_CreateDriver_Success(t *testing.T) {
	ctrl, driverRepo, userRepo, authRepo := setupDriverController(t)

	userEntity := &entities.PublicUser{
		User:  entities.User{ID: "user-1", RefID: 1, AuthRefID: 10},
		Email: "test@example.com",
	}
	driverEntity := &entities.Driver{ID: "drv-1", RefID: 1, DriverLicense: "DL-12345"}

	userRepo.EXPECT().FindByID(mock.Anything, "user-1").Return(userEntity, nil)
	driverRepo.EXPECT().FindByUserRefID(mock.Anything, int64(1)).Return(nil, nil)
	driverRepo.EXPECT().Create(mock.Anything, entities.CreateDriverData{
		DriverLicense: "DL-12345",
		UserRefID:     1,
	}).Return(driverEntity, nil)
	authRepo.EXPECT().UpdateRole(mock.Anything, int64(10), "DRIVER").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/drivers", ctrl.CreateDriver)

	body := `{"driverLicense":"DL-12345"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/drivers", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestDriverController_CreateDriver_InvalidJSON(t *testing.T) {
	ctrl, _, _, _ := setupDriverController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/drivers", ctrl.CreateDriver)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/drivers", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDriverController_CreateDriver_ValidationError(t *testing.T) {
	ctrl, _, _, _ := setupDriverController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/drivers", ctrl.CreateDriver)

	body := `{"driverLicense":""}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/drivers", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDriverController_CreateDriver_AlreadyExists(t *testing.T) {
	ctrl, driverRepo, userRepo, _ := setupDriverController(t)

	userEntity := &entities.PublicUser{
		User:  entities.User{ID: "user-1", RefID: 1, AuthRefID: 10},
		Email: "test@example.com",
	}
	existingDriver := &entities.Driver{ID: "drv-1", RefID: 1}

	userRepo.EXPECT().FindByID(mock.Anything, "user-1").Return(userEntity, nil)
	driverRepo.EXPECT().FindByUserRefID(mock.Anything, int64(1)).Return(existingDriver, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/drivers", ctrl.CreateDriver)

	body := `{"driverLicense":"DL-12345"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/drivers", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// c.Error() is called with DriverAlreadyExistsError
	assert.Equal(t, http.StatusOK, w.Code)
}
