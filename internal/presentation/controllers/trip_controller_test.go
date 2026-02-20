package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/trip"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTripController(t *testing.T) (
	*TripController,
	*mocks.MockTripRepository,
	*mocks.MockDriverRepository,
	*mocks.MockCarRepository,
	*mocks.MockCityRepository,
) {
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)
	carRepo := mocks.NewMockCarRepository(t)
	cityRepo := mocks.NewMockCityRepository(t)

	listUC := trip.NewListTripsUseCase(tripRepo)
	getUC := trip.NewGetTripUseCase(tripRepo)
	findUC := trip.NewFindTripsUseCase(tripRepo)
	createUC := trip.NewCreateTripUseCase(tripRepo, driverRepo, carRepo, cityRepo)
	deleteUC := trip.NewDeleteTripUseCase(tripRepo, driverRepo)
	ctrl := NewTripController(listUC, getUC, findUC, createUC, deleteUC)

	return ctrl, tripRepo, driverRepo, carRepo, cityRepo
}

func TestTripController_ListTrips_Success(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	trips := []entities.Trip{
		{ID: "trip-1", Kms: 100, Seats: 3},
	}
	tripRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(trips, 1, nil)

	router := gin.New()
	router.GET("/trips", ctrl.ListTrips)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestTripController_ListTrips_Empty(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	tripRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return([]entities.Trip{}, 0, nil)

	router := gin.New()
	router.GET("/trips", ctrl.ListTrips)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTripController_ListTrips_Error(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	tripRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/trips", ctrl.ListTrips)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTripController_GetTrip_Success(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	tripEntity := &entities.Trip{ID: "trip-1", Kms: 200, Seats: 4}
	tripRepo.EXPECT().FindByID(mock.Anything, "trip-1").Return(tripEntity, nil)

	router := gin.New()
	router.GET("/trips/:id", ctrl.GetTrip)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips/trip-1", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestTripController_GetTrip_NotFound(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	tripRepo.EXPECT().FindByID(mock.Anything, "trip-999").Return(nil, nil)

	router := gin.New()
	router.GET("/trips/:id", ctrl.GetTrip)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips/trip-999", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTripController_FindTrip_Success(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	trips := []entities.Trip{
		{ID: "trip-1", Kms: 200},
	}
	tripRepo.EXPECT().FindByFilters(mock.Anything, mock.Anything).Return(trips, nil)

	router := gin.New()
	router.GET("/trips/search", ctrl.FindTrip)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips/search?departureCity=Paris&arrivalCity=Lyon", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestTripController_FindTrip_WithDate(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	trips := []entities.Trip{}
	tripRepo.EXPECT().FindByFilters(mock.Anything, mock.Anything).Return(trips, nil)

	router := gin.New()
	router.GET("/trips/search", ctrl.FindTrip)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips/search?date=2025-06-15", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTripController_CreateTrip_InvalidJSON(t *testing.T) {
	ctrl, _, _, _, _ := setupTripController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/trips", ctrl.CreateTrip)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/trips", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTripController_CreateTrip_ValidationError(t *testing.T) {
	ctrl, _, _, _, _ := setupTripController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/trips", ctrl.CreateTrip)

	body := `{"kms":0,"date":"","departureCity":"","arrivalCity":"","seats":0,"carId":""}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/trips", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTripController_DeleteTrip_Success(t *testing.T) {
	ctrl, tripRepo, driverRepo, _, _ := setupTripController(t)

	existing := &entities.Trip{ID: "trip-1", DriverRefID: 1}
	driver := &entities.Driver{ID: "drv-1", RefID: 1}

	tripRepo.EXPECT().FindByID(mock.Anything, "trip-1").Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, "user-1").Return(driver, nil)
	tripRepo.EXPECT().Delete(mock.Anything, "trip-1").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.DELETE("/trips/:id", ctrl.DeleteTrip)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/trips/trip-1", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTripController_DeleteTrip_NotFound(t *testing.T) {
	ctrl, tripRepo, _, _, _ := setupTripController(t)

	tripRepo.EXPECT().FindByID(mock.Anything, "trip-999").Return(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.DELETE("/trips/:id", ctrl.DeleteTrip)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/trips/trip-999", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
