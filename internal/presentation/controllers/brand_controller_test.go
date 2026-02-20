package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/brand"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupBrandController(t *testing.T) (*BrandController, *mocks.MockBrandRepository) {
	brandRepo := mocks.NewMockBrandRepository(t)

	listUC := brand.NewListBrandsUseCase(brandRepo)
	createUC := brand.NewCreateBrandUseCase(brandRepo)
	deleteUC := brand.NewDeleteBrandUseCase(brandRepo)
	ctrl := NewBrandController(listUC, createUC, deleteUC)

	return ctrl, brandRepo
}

func TestBrandController_ListBrands_Success(t *testing.T) {
	ctrl, brandRepo := setupBrandController(t)

	brands := []entities.Brand{
		{ID: "brand-1", Name: "Toyota"},
		{ID: "brand-2", Name: "Honda"},
	}
	brandRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(brands, 2, nil)

	router := gin.New()
	router.GET("/brands", ctrl.ListBrands)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/brands", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
	assert.NotNil(t, resp["data"])
	assert.NotNil(t, resp["meta"])
}

func TestBrandController_ListBrands_WithPagination(t *testing.T) {
	ctrl, brandRepo := setupBrandController(t)

	brands := []entities.Brand{{ID: "brand-3", Name: "BMW"}}
	brandRepo.EXPECT().FindAll(mock.Anything, 10, 10).Return(brands, 11, nil)

	router := gin.New()
	router.GET("/brands", ctrl.ListBrands)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/brands?page=2&limit=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandController_ListBrands_Error(t *testing.T) {
	ctrl, brandRepo := setupBrandController(t)

	brandRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("db error"))

	router := gin.New()
	router.GET("/brands", ctrl.ListBrands)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/brands", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandController_CreateBrand_Success(t *testing.T) {
	ctrl, brandRepo := setupBrandController(t)

	newBrand := &entities.Brand{ID: "brand-1", Name: "Tesla"}
	brandRepo.EXPECT().Create(mock.Anything, entities.CreateBrandData{Name: "Tesla"}).Return(newBrand, nil)

	router := gin.New()
	router.POST("/brands", ctrl.CreateBrand)

	body := `{"name":"Tesla"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/brands", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestBrandController_CreateBrand_InvalidJSON(t *testing.T) {
	ctrl, _ := setupBrandController(t)

	router := gin.New()
	router.POST("/brands", ctrl.CreateBrand)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/brands", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBrandController_CreateBrand_ValidationError(t *testing.T) {
	ctrl, _ := setupBrandController(t)

	router := gin.New()
	router.POST("/brands", ctrl.CreateBrand)

	body := `{"name":""}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/brands", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBrandController_DeleteBrand_Success(t *testing.T) {
	ctrl, brandRepo := setupBrandController(t)

	existing := &entities.Brand{ID: "brand-1", Name: "Tesla"}
	brandRepo.EXPECT().FindByID(mock.Anything, "brand-1").Return(existing, nil)
	brandRepo.EXPECT().Delete(mock.Anything, "brand-1").Return(nil)

	router := gin.New()
	router.DELETE("/brands/:id", ctrl.DeleteBrand)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/brands/brand-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestBrandController_DeleteBrand_NotFound(t *testing.T) {
	ctrl, brandRepo := setupBrandController(t)

	brandRepo.EXPECT().FindByID(mock.Anything, "brand-999").Return(nil, nil)

	router := gin.New()
	router.DELETE("/brands/:id", ctrl.DeleteBrand)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/brands/brand-999", nil)
	router.ServeHTTP(w, req)

	// c.Error() is called with BrandNotFoundError
	assert.Equal(t, http.StatusOK, w.Code)
}
