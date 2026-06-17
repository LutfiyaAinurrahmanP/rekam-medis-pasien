package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	userservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupUserRouter() (*gin.Engine, *routes.RouteConfig) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	userRepo := repository.NewUserRepository(db)
	userService := userservice.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService)

	routeCfg := &routes.RouteConfig{
		Config:      cfg,
		UserHandler: userHandler,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupUsersRouter(v1, routeCfg)

	// create super admin for hard-delete test
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	sa := &models.User{
		Username: "superadmin_test",
		Email:    "sa@test.com",
		Phone:    "08999999999",
		Password: string(password),
		Role:     models.RoleSuperAdmin,
		IsActive: true,
	}
	userRepo.Create(sa)

	return r, routeCfg
}

// Helper to create a user and get its ID and token
func createTestUserAndGetToken(r *gin.Engine, cfg *routes.RouteConfig, role string, username string) (uint, string) {
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	reqBody := dto.CreateUserRequest{
		Username: username,
		Email:    username + "@example.com",
		Phone:    "081234567890",
		Password: "password123",
		Role:     role,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	userID := res["data"].(map[string]interface{})["id"].(float64)

	token := GenerateTestToken(uint(userID), role, cfg.Config)
	return uint(userID), token
}

func TestIntegration_User_CreateUser(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreateUserRequest{
		Username: "newuser123",
		Email:    "newuser@example.com",
		Phone:    "081234567890",
		Password: "password123",
		Role:     models.RolePatient,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_User_CreateUser_ForbiddenForPatient(t *testing.T) {
	r, cfg := setupUserRouter()
	patientToken := GenerateTestToken(1, models.RolePatient, cfg.Config)

	reqBody := dto.CreateUserRequest{
		Username: "newuser123",
		Email:    "newuser@example.com",
		Phone:    "081234567890",
		Password: "password123",
		Role:     models.RolePatient,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+patientToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------
// /users/me Routes
// ---------------------------------------------------------

func TestIntegration_User_GetMyProfile(t *testing.T) {
	r, cfg := setupUserRouter()
	_, userToken := createTestUserAndGetToken(r, cfg, models.RolePatient, "profileuser")

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_UpdateMyProfile(t *testing.T) {
	r, cfg := setupUserRouter()
	_, userToken := createTestUserAndGetToken(r, cfg, models.RolePatient, "updateme")

	newName := "updatedname"
	updateReq := dto.UpdateUserRequest{
		Username: &newName,
	}
	updateBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_ChangeMyPassword(t *testing.T) {
	r, cfg := setupUserRouter()
	_, userToken := createTestUserAndGetToken(r, cfg, models.RolePatient, "changepw")

	reqBody := dto.ChangePasswordRequest{
		OldPassword: "password123",
		NewPassword: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/users/me/change-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_DeleteMyAccount(t *testing.T) {
	r, cfg := setupUserRouter()
	_, userToken := createTestUserAndGetToken(r, cfg, models.RolePatient, "deleteme")

	reqBody := dto.DeleteAccountRequest{
		Password: "password123",
		Reason:   "No longer need the account",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_DeactivateMyAccount(t *testing.T) {
	r, cfg := setupUserRouter()
	_, userToken := createTestUserAndGetToken(r, cfg, models.RolePatient, "deactivateme")

	reqBody := dto.DeactivateAccountRequest{
		Password: "password123",
		Reason:   "Taking a break",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/users/me/deactivate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ---------------------------------------------------------
// Admin Routes (/users)
// ---------------------------------------------------------

func TestIntegration_User_ListUsers(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_DeleteListUsers(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/deleted", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_GetUserByID(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "getme")

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", userID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_UpdateUser(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "updateuser")

	newName := "updatednameadmin"
	updateReq := dto.UpdateUserRequest{
		Username: &newName,
	}
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_SoftDeleteUser(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "softdelete")

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", userID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_RestoreUser(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "restoreuser")

	// Delete first
	reqDel, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", userID), nil)
	reqDel.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(httptest.NewRecorder(), reqDel)

	// Restore
	reqRes, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%d/restore", userID), nil)
	reqRes.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqRes)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_ResetPassword(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "resetpw")

	reqBody := dto.ResetPasswordRequest{
		NewPassword: "newpassword456",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%d/reset-password", userID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_DeactivateAndActivateUser(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "toggleactive")

	// Deactivate
	reqDeac, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%d/deactivate", userID), nil)
	reqDeac.Header.Set("Authorization", "Bearer "+adminToken)
	wDeac := httptest.NewRecorder()
	r.ServeHTTP(wDeac, reqDeac)
	assert.Equal(t, http.StatusOK, wDeac.Code)

	// Activate
	reqAc, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/users/%d/activate", userID), nil)
	reqAc.Header.Set("Authorization", "Bearer "+adminToken)
	wAc := httptest.NewRecorder()
	r.ServeHTTP(wAc, reqAc)
	assert.Equal(t, http.StatusOK, wAc.Code)
}

// ---------------------------------------------------------
// SuperAdmin Routes
// ---------------------------------------------------------

func TestIntegration_User_HardDeleteUser(t *testing.T) {
	r, cfg := setupUserRouter()
	saToken := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "harddelete")

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d/hard-delete", userID), nil)
	req.Header.Set("Authorization", "Bearer "+saToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_User_HardDeleteUser_Forbidden(t *testing.T) {
	r, cfg := setupUserRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	userID, _ := createTestUserAndGetToken(r, cfg, models.RolePatient, "harddelete_forbidden")

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d/hard-delete", userID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Admin is not SuperAdmin, so it should be forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
}
