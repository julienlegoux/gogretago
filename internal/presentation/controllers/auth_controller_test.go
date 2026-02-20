package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/application/usecases/auth"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/services"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type authControllerDeps struct {
	ctrl        *AuthController
	authRepo    *mocks.MockAuthRepository
	userRepo    *mocks.MockUserRepository
	passwordSvc *mocks.MockPasswordService
	emailSvc    *mocks.MockEmailService
	jwtSvc      *mocks.MockJwtService
}

func setupAuthController(t *testing.T) authControllerDeps {
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	registerUC := auth.NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	loginUC := auth.NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	ctrl := NewAuthController(registerUC, loginUC)

	return authControllerDeps{ctrl, authRepo, userRepo, passwordSvc, emailSvc, jwtSvc}
}

func TestAuthController_Register_Success(t *testing.T) {
	d := setupAuthController(t)
	ctrl, authRepo, passwordSvc, emailSvc, jwtSvc := d.ctrl, d.authRepo, d.passwordSvc, d.emailSvc, d.jwtSvc

	authEntity := &entities.Auth{ID: "auth-1", RefID: 1, Email: "test@example.com", Role: "USER"}
	userEntity := &entities.PublicUser{
		User:  entities.User{ID: "user-1", RefID: 1, AuthRefID: 1},
		Email: "test@example.com",
	}

	authRepo.EXPECT().ExistsByEmail(mock.Anything, "test@example.com").Return(false, nil)
	passwordSvc.EXPECT().Hash("Password1").Return("hashed", nil)
	authRepo.EXPECT().CreateWithUser(mock.Anything,
		entities.CreateAuthData{Email: "test@example.com", Password: "hashed"},
		entities.CreateUserData{},
	).Return(authEntity, userEntity, nil)
	emailSvc.EXPECT().SendWelcomeEmail("test@example.com", "there").Return(nil)
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("jwt-token", nil)

	router := gin.New()
	router.POST("/register", ctrl.Register)

	body := `{"email":"test@example.com","password":"Password1","confirmPassword":"Password1"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestAuthController_Register_InvalidJSON(t *testing.T) {
	d := setupAuthController(t)
	ctrl := d.ctrl

	router := gin.New()
	router.POST("/register", ctrl.Register)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["success"])
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "VALIDATION_ERROR", errObj["code"])
}

func TestAuthController_Register_ValidationError(t *testing.T) {
	d := setupAuthController(t)
	ctrl := d.ctrl

	router := gin.New()
	router.POST("/register", ctrl.Register)

	body := `{"email":"test@example.com"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["success"])
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "VALIDATION_ERROR", errObj["code"])
	assert.Equal(t, "Validation failed", errObj["message"])
	assert.NotNil(t, errObj["details"])
}

func TestAuthController_Register_UseCaseError(t *testing.T) {
	d := setupAuthController(t)
	ctrl, authRepo := d.ctrl, d.authRepo

	authRepo.EXPECT().ExistsByEmail(mock.Anything, "existing@example.com").Return(true, nil)

	router := gin.New()
	router.POST("/register", ctrl.Register)

	body := `{"email":"existing@example.com","password":"Password1","confirmPassword":"Password1"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// The controller calls c.Error() which doesn't set status by itself
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthController_Login_Success(t *testing.T) {
	d := setupAuthController(t)
	ctrl, authRepo, userRepo, passwordSvc, jwtSvc := d.ctrl, d.authRepo, d.userRepo, d.passwordSvc, d.jwtSvc

	authEntity := &entities.Auth{ID: "auth-1", RefID: 1, Email: "test@example.com", Password: "hashed", Role: "USER"}
	userEntity := &entities.PublicUser{
		User:  entities.User{ID: "user-1", RefID: 1, AuthRefID: 1},
		Email: "test@example.com",
	}

	authRepo.EXPECT().FindByEmail(mock.Anything, "test@example.com").Return(authEntity, nil)
	passwordSvc.EXPECT().Verify("Password1", "hashed").Return(true, nil)
	userRepo.EXPECT().FindByAuthRefID(mock.Anything, int64(1)).Return(userEntity, nil)
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("jwt-token", nil)

	router := gin.New()
	router.POST("/login", ctrl.Login)

	body := `{"email":"test@example.com","password":"Password1"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["success"])
}

func TestAuthController_Login_InvalidJSON(t *testing.T) {
	d := setupAuthController(t)
	ctrl := d.ctrl

	router := gin.New()
	router.POST("/login", ctrl.Login)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`not json`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_Login_ValidationError(t *testing.T) {
	d := setupAuthController(t)
	ctrl := d.ctrl

	router := gin.New()
	router.POST("/login", ctrl.Login)

	body := `{"email":"not-an-email","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	errObj := resp["error"].(map[string]interface{})
	assert.Equal(t, "VALIDATION_ERROR", errObj["code"])
}

func TestAuthController_Login_UseCaseError(t *testing.T) {
	d := setupAuthController(t)
	ctrl, authRepo := d.ctrl, d.authRepo

	authRepo.EXPECT().FindByEmail(mock.Anything, "bad@example.com").Return(nil, fmt.Errorf("db error"))

	router := gin.New()
	router.POST("/login", ctrl.Login)

	body := `{"email":"bad@example.com","password":"Password1"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
