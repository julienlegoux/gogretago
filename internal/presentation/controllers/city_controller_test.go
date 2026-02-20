package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/city"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupCityController(t *testing.T) (*CityController, *mocks.MockCityRepository) {
	cityRepo := mocks.NewMockCityRepository(t)

	listUC := city.NewListCitiesUseCase(cityRepo)
	createUC := city.NewCreateCityUseCase(cityRepo)
	deleteUC := city.NewDeleteCityUseCase(cityRepo)
	ctrl := NewCityController(listUC, createUC, deleteUC)

	return ctrl, cityRepo
}

func TestCityController_ListCities_Success(t *testing.T) {
	ctrl, cityRepo := setupCityController(t)

	cities := []entities.City{
		{ID: "city-1", CityName: "Paris", Zipcode: "75000"},
	}
	cityRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(cities, 1, nil)

	router := gin.New()
	router.GET("/cities", ctrl.ListCities)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cities", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestCityController_ListCities_Empty(t *testing.T) {
	ctrl, cityRepo := setupCityController(t)

	cityRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return([]entities.City{}, 0, nil)

	router := gin.New()
	router.GET("/cities", ctrl.ListCities)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cities", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCityController_ListCities_Error(t *testing.T) {
	ctrl, cityRepo := setupCityController(t)

	cityRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/cities", ctrl.ListCities)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cities", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCityController_CreateCity_Success(t *testing.T) {
	ctrl, cityRepo := setupCityController(t)

	newCity := &entities.City{ID: "city-1", CityName: "Lyon", Zipcode: "69000"}
	cityRepo.EXPECT().Create(mock.Anything, entities.CreateCityData{CityName: "Lyon", Zipcode: "69000"}).Return(newCity, nil)

	router := gin.New()
	router.POST("/cities", ctrl.CreateCity)

	body := `{"cityName":"Lyon","zipcode":"69000"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/cities", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCityController_CreateCity_InvalidJSON(t *testing.T) {
	ctrl, _ := setupCityController(t)

	router := gin.New()
	router.POST("/cities", ctrl.CreateCity)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/cities", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCityController_CreateCity_ValidationError(t *testing.T) {
	ctrl, _ := setupCityController(t)

	router := gin.New()
	router.POST("/cities", ctrl.CreateCity)

	body := `{"cityName":"","zipcode":""}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/cities", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCityController_DeleteCity_Success(t *testing.T) {
	ctrl, cityRepo := setupCityController(t)

	existing := &entities.City{ID: "city-1", CityName: "Paris"}
	cityRepo.EXPECT().FindByID(mock.Anything, "city-1").Return(existing, nil)
	cityRepo.EXPECT().Delete(mock.Anything, "city-1").Return(nil)

	router := gin.New()
	router.DELETE("/cities/:id", ctrl.DeleteCity)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/cities/city-1", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCityController_DeleteCity_NotFound(t *testing.T) {
	ctrl, cityRepo := setupCityController(t)

	cityRepo.EXPECT().FindByID(mock.Anything, "city-999").Return(nil, nil)

	router := gin.New()
	router.DELETE("/cities/:id", ctrl.DeleteCity)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/cities/city-999", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
