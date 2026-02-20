package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/inscription"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupInscriptionController(t *testing.T) (
	*InscriptionController,
	*mocks.MockInscriptionRepository,
	*mocks.MockUserRepository,
	*mocks.MockTripRepository,
) {
	inscRepo := mocks.NewMockInscriptionRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	tripRepo := mocks.NewMockTripRepository(t)

	listUC := inscription.NewListInscriptionsUseCase(inscRepo)
	createUC := inscription.NewCreateInscriptionUseCase(inscRepo, userRepo, tripRepo)
	deleteUC := inscription.NewDeleteInscriptionUseCase(inscRepo)
	listUserUC := inscription.NewListUserInscriptionsUseCase(inscRepo)
	listPassengersUC := inscription.NewListTripPassengersUseCase(inscRepo)
	ctrl := NewInscriptionController(listUC, createUC, deleteUC, listUserUC, listPassengersUC)

	return ctrl, inscRepo, userRepo, tripRepo
}

func TestInscriptionController_ListInscriptions_Success(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscriptions := []entities.Inscription{
		{ID: "insc-1", UserRefID: 1, TripRefID: 1},
	}
	inscRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(inscriptions, 1, nil)

	router := gin.New()
	router.GET("/inscriptions", ctrl.ListInscriptions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/inscriptions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestInscriptionController_ListInscriptions_Empty(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return([]entities.Inscription{}, 0, nil)

	router := gin.New()
	router.GET("/inscriptions", ctrl.ListInscriptions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/inscriptions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInscriptionController_ListInscriptions_Error(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/inscriptions", ctrl.ListInscriptions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/inscriptions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInscriptionController_CreateInscription_Success(t *testing.T) {
	ctrl, inscRepo, userRepo, tripRepo := setupInscriptionController(t)

	userEntity := &entities.PublicUser{
		User:  entities.User{ID: "user-1", RefID: 1},
		Email: "test@example.com",
	}
	tripEntity := &entities.Trip{ID: "trip-1", RefID: 10, Seats: 3}
	newInsc := &entities.Inscription{ID: "insc-1", UserRefID: 1, TripRefID: 10}

	userRepo.EXPECT().FindByID(mock.Anything, "user-1").Return(userEntity, nil)
	tripRepo.EXPECT().FindByID(mock.Anything, "trip-1").Return(tripEntity, nil)
	inscRepo.EXPECT().ExistsByUserAndTrip(mock.Anything, int64(1), int64(10)).Return(false, nil)
	inscRepo.EXPECT().CountByTripRefID(mock.Anything, int64(10)).Return(0, nil)
	inscRepo.EXPECT().Create(mock.Anything, entities.CreateInscriptionData{
		UserRefID: 1,
		TripRefID: 10,
	}).Return(newInsc, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/inscriptions", ctrl.CreateInscription)

	body := `{"tripId":"trip-1"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inscriptions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestInscriptionController_CreateInscription_InvalidJSON(t *testing.T) {
	ctrl, _, _, _ := setupInscriptionController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/inscriptions", ctrl.CreateInscription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inscriptions", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInscriptionController_CreateInscription_ValidationError(t *testing.T) {
	ctrl, _, _, _ := setupInscriptionController(t)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.POST("/inscriptions", ctrl.CreateInscription)

	body := `{"tripId":""}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inscriptions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInscriptionController_DeleteInscription_Success(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	existing := &entities.Inscription{ID: "insc-1", UserRefID: 1}
	inscRepo.EXPECT().FindByIDAndUserID(mock.Anything, "insc-1", "user-1").Return(existing, nil)
	inscRepo.EXPECT().Delete(mock.Anything, "insc-1").Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.DELETE("/inscriptions/:id", ctrl.DeleteInscription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/inscriptions/insc-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestInscriptionController_DeleteInscription_NotFound(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscRepo.EXPECT().FindByIDAndUserID(mock.Anything, "insc-999", "user-1").Return(nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Next()
	})
	router.DELETE("/inscriptions/:id", ctrl.DeleteInscription)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/inscriptions/insc-999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInscriptionController_ListUserInscriptions_Success(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscriptions := []entities.Inscription{
		{ID: "insc-1", UserRefID: 1},
	}
	inscRepo.EXPECT().FindByUserID(mock.Anything, "user-1").Return(inscriptions, nil)

	router := gin.New()
	router.GET("/users/:id/inscriptions", ctrl.ListUserInscriptions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/user-1/inscriptions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestInscriptionController_ListUserInscriptions_Error(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscRepo.EXPECT().FindByUserID(mock.Anything, "user-1").Return(nil, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/users/:id/inscriptions", ctrl.ListUserInscriptions)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/user-1/inscriptions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInscriptionController_ListTripPassengers_Success(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscriptions := []entities.Inscription{
		{ID: "insc-1", TripRefID: 1},
	}
	inscRepo.EXPECT().FindByTripID(mock.Anything, "trip-1").Return(inscriptions, nil)

	router := gin.New()
	router.GET("/trips/:id/passengers", ctrl.ListTripPassengers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips/trip-1/passengers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestInscriptionController_ListTripPassengers_Error(t *testing.T) {
	ctrl, inscRepo, _, _ := setupInscriptionController(t)

	inscRepo.EXPECT().FindByTripID(mock.Anything, "trip-1").Return(nil, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/trips/:id/passengers", ctrl.ListTripPassengers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trips/trip-1/passengers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
