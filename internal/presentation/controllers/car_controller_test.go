package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/car"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupCarController(t *testing.T) (
	*CarController,
	*mocks.MockCarRepository,
	*mocks.MockModelRepository,
	*mocks.MockBrandRepository,
	*mocks.MockDriverRepository,
) {
	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	listUC := car.NewListCarsUseCase(carRepo)
	createUC := car.NewCreateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	updateUC := car.NewUpdateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	deleteUC := car.NewDeleteCarUseCase(carRepo, driverRepo)
	ctrl := NewCarController(listUC, createUC, updateUC, deleteUC)

	return ctrl, carRepo, modelRepo, brandRepo, driverRepo
}

func TestCarController_ListCars_Success(t *testing.T) {
	ctrl, carRepo, _, _, _ := setupCarController(t)

	cars := []entities.Car{
		{ID: "car-1", LicensePlate: "AB-123-CD"},
	}
	carRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(cars, 1, nil)

	router := gin.New()
	router.GET("/cars", ctrl.ListCars)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cars", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestCarController_ListCars_Empty(t *testing.T) {
	ctrl, carRepo, _, _, _ := setupCarController(t)

	carRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return([]entities.Car{}, 0, nil)

	router := gin.New()
	router.GET("/cars", ctrl.ListCars)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cars", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCarController_ListCars_Error(t *testing.T) {
	ctrl, carRepo, _, _, _ := setupCarController(t)

	carRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/cars", ctrl.ListCars)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cars", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCarController_CreateCar_Success(t *testing.T) {
	ctrl, carRepo, modelRepo, brandRepo, driverRepo := setupCarController(t)

	driver := &entities.Driver{ID: "drv-1", RefID: 1}
	brand := &entities.Brand{ID: "brand-1", RefID: 10, Name: "Toyota"}
	model := &entities.VehicleModel{ID: "model-1", RefID: 100, Name: "Corolla"}
	newCar := &entities.Car{ID: "car-1", LicensePlate: "AB-123-CD"}

	driverRepo.EXPECT().FindByUserID(mock.Anything, "user-1").Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, "AB-123-CD").Return(false, nil)
	brandRepo.EXPECT().FindByID(mock.Anything, "brand-1").Return(brand, nil)
	modelRepo.EXPECT().FindByNameAndBrand(mock.Anything, "Corolla", int64(10)).Return(model, nil)
	carRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(newCar, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/cars", ctrl.CreateCar)

	body := `{"model":"Corolla","brandId":"brand-1","licensePlate":"AB-123-CD"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/cars", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCarController_CreateCar_InvalidJSON(t *testing.T) {
	ctrl, _, _, _, _ := setupCarController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/cars", ctrl.CreateCar)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/cars", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarController_CreateCar_ValidationError(t *testing.T) {
	ctrl, _, _, _, _ := setupCarController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/cars", ctrl.CreateCar)

	body := `{"model":"","brandId":"","licensePlate":""}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/cars", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarController_UpdateCar_InvalidJSON(t *testing.T) {
	ctrl, _, _, _, _ := setupCarController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.PUT("/cars/:id", ctrl.UpdateCar)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/cars/car-1", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarController_PatchCar_InvalidJSON(t *testing.T) {
	ctrl, _, _, _, _ := setupCarController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.PATCH("/cars/:id", ctrl.PatchCar)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/cars/car-1", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarController_DeleteCar_Success(t *testing.T) {
	ctrl, carRepo, _, _, driverRepo := setupCarController(t)

	existing := &entities.Car{ID: "car-1", DriverRefID: 1}
	driver := &entities.Driver{ID: "drv-1", RefID: 1}

	carRepo.EXPECT().FindByID(mock.Anything, "car-1").Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, "user-1").Return(driver, nil)
	carRepo.EXPECT().Delete(mock.Anything, "car-1").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.DELETE("/cars/:id", ctrl.DeleteCar)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/cars/car-1", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCarController_DeleteCar_NotFound(t *testing.T) {
	ctrl, carRepo, _, _, _ := setupCarController(t)

	carRepo.EXPECT().FindByID(mock.Anything, "car-999").Return(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.DELETE("/cars/:id", ctrl.DeleteCar)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/cars/car-999", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
