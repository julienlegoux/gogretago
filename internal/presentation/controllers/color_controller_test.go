package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/color"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupColorController(t *testing.T) (*ColorController, *mocks.MockColorRepository) {
	colorRepo := mocks.NewMockColorRepository(t)

	listUC := color.NewListColorsUseCase(colorRepo)
	createUC := color.NewCreateColorUseCase(colorRepo)
	updateUC := color.NewUpdateColorUseCase(colorRepo)
	deleteUC := color.NewDeleteColorUseCase(colorRepo)
	ctrl := NewColorController(listUC, createUC, updateUC, deleteUC)

	return ctrl, colorRepo
}

func TestColorController_ListColors_Success(t *testing.T) {
	ctrl, colorRepo := setupColorController(t)

	colors := []entities.Color{
		{ID: "color-1", Name: "Red", Hex: "#FF0000"},
	}
	colorRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(colors, 1, nil)

	router := gin.New()
	router.GET("/colors", ctrl.ListColors)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/colors", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestColorController_ListColors_Empty(t *testing.T) {
	ctrl, colorRepo := setupColorController(t)

	colorRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return([]entities.Color{}, 0, nil)

	router := gin.New()
	router.GET("/colors", ctrl.ListColors)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/colors", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestColorController_ListColors_Error(t *testing.T) {
	ctrl, colorRepo := setupColorController(t)

	colorRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/colors", ctrl.ListColors)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/colors", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestColorController_CreateColor_Success(t *testing.T) {
	ctrl, colorRepo := setupColorController(t)

	colorRepo.EXPECT().FindByName(mock.Anything, "Blue").Return(nil, nil)
	newColor := &entities.Color{ID: "color-1", Name: "Blue", Hex: "#0000FF"}
	colorRepo.EXPECT().Create(mock.Anything, entities.CreateColorData{Name: "Blue", Hex: "#0000FF"}).Return(newColor, nil)

	router := gin.New()
	router.POST("/colors", ctrl.CreateColor)

	body := `{"name":"Blue","hex":"#0000FF"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/colors", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestColorController_CreateColor_InvalidJSON(t *testing.T) {
	ctrl, _ := setupColorController(t)

	router := gin.New()
	router.POST("/colors", ctrl.CreateColor)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/colors", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestColorController_CreateColor_ValidationError(t *testing.T) {
	ctrl, _ := setupColorController(t)

	router := gin.New()
	router.POST("/colors", ctrl.CreateColor)

	body := `{"name":"","hex":"invalid"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/colors", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestColorController_UpdateColor_Success(t *testing.T) {
	ctrl, colorRepo := setupColorController(t)

	existing := &entities.Color{ID: "color-1", Name: "Red", Hex: "#FF0000"}
	colorRepo.EXPECT().FindByID(mock.Anything, "color-1").Return(existing, nil)
	newName := "Crimson"
	colorRepo.EXPECT().FindByName(mock.Anything, "Crimson").Return(nil, nil)
	updated := &entities.Color{ID: "color-1", Name: "Crimson", Hex: "#FF0000"}
	colorRepo.EXPECT().Update(mock.Anything, "color-1", mock.Anything).Return(updated, nil)

	router := gin.New()
	router.PATCH("/colors/:id", ctrl.UpdateColor)

	body := `{"name":"` + newName + `"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/colors/color-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestColorController_UpdateColor_InvalidJSON(t *testing.T) {
	ctrl, _ := setupColorController(t)

	router := gin.New()
	router.PATCH("/colors/:id", ctrl.UpdateColor)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/colors/color-1", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestColorController_DeleteColor_Success(t *testing.T) {
	ctrl, colorRepo := setupColorController(t)

	existing := &entities.Color{ID: "color-1", Name: "Red"}
	colorRepo.EXPECT().FindByID(mock.Anything, "color-1").Return(existing, nil)
	colorRepo.EXPECT().Delete(mock.Anything, "color-1").Return(nil)

	router := gin.New()
	router.DELETE("/colors/:id", ctrl.DeleteColor)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/colors/color-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestColorController_DeleteColor_NotFound(t *testing.T) {
	ctrl, colorRepo := setupColorController(t)

	colorRepo.EXPECT().FindByID(mock.Anything, "color-999").Return(nil, nil)

	router := gin.New()
	router.DELETE("/colors/:id", ctrl.DeleteColor)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/colors/color-999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
