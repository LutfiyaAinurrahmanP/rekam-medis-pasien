package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/database"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	userservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/user"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	realTestDB     *gorm.DB
	realTestRouter *gin.Engine
	realTestConfig *config.Config
)

// Setup test router dengan database real
func setupRealTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Initialize config
	realTestConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-real-testing",
			ExpiredTime: 24 * 60 * time.Minute,
		},
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	// Initialize database (SQLite in-memory)
	realTestDB = database.InitTestDB()

	// Run migrations
	realTestDB.AutoMigrate(&models.User{})

	// Initialize repositories
	userRepo := repository.NewUserRepository(realTestDB)

	// Initialize services
	userService := userservice.NewUserService(userRepo, realTestConfig)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := gin.New()

	// Setup routes
	api := router.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Config:      realTestConfig,
		UserHandler: userHandler,
	}
	routes.SetupAPIRouter(api, routeConfig)

	return router
}

// Cleanup database sebelum test
func cleanupRealTestDB() {
	realTestDB.Exec("DELETE FROM users")
}

// Helper function untuk membuat user langsung di database
func createRealUser(username, email, phone, password, role string, isActive bool) (*models.User, string, error) {
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: isActive,
	}

	result := realTestDB.Create(user)
	if result.Error != nil {
		return nil, "", result.Error
	}

	token, _, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role,
		realTestConfig.JWT.Secret, realTestConfig.JWT.ExpiredTime)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Helper function untuk perform request
func performRealRequest(method, path, token string, bodyJSON map[string]interface{}) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if bodyJSON != nil {
		bodyBytes, _ = json.Marshal(bodyJSON)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	realTestRouter.ServeHTTP(w, req)
	return w
}

func init() {
	realTestRouter = setupRealTestRouter()
}

// ==================== GET /api/v1/users/me Tests ====================

func Test_GetMyProfile_Success(t *testing.T) {
	cleanupRealTestDB()

	// Buat user terlebih dahulu
	user, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	// Request GET /api/v1/users/me
	w := performRealRequest("GET", "/api/v1/users/me", token, nil)

	// Validasi response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Logf("JSON Unmarshal Error: %v", err)
		t.FailNow()
	}

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, user.Username, data["username"])
	assert.Equal(t, user.Email, data["email"])
	assert.Equal(t, user.Phone, data["phone"])
	assert.Equal(t, user.Role, data["role"])
	assert.Equal(t, user.IsActive, data["is_active"])

	t.Logf("[PASS] GET /api/v1/users/me - Success")
}

func Test_GetMyProfile_Unauthorized(t *testing.T) {
	cleanupRealTestDB()

	// Request tanpa token
	w := performRealRequest("GET", "/api/v1/users/me", "", nil)

	// Validasi response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/users/me - Unauthorized")
}

func Test_GetMyProfile_InvalidToken(t *testing.T) {
	cleanupRealTestDB()

	// Request dengan token invalid
	w := performRealRequest("GET", "/api/v1/users/me", "invalid-token-12345", nil)

	// Validasi response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/users/me - Invalid Token")
}

// ==================== PUT /api/v1/users/me Tests ====================

func Test_UpdateMyProfile_Success(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("oldusername", "old@example.com", "081111111111", "password123", models.RolePatient, true)

	// Request body seperti di Postman
	requestBody := map[string]interface{}{
		"username": "newusername",
		"email":    "new@example.com",
		"phone":    "081222222222",
	}

	// Request PUT /api/v1/users/me
	w := performRealRequest("PUT", "/api/v1/users/me", token, requestBody)

	// Validasi response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "newusername", data["username"])
	assert.Equal(t, "new@example.com", data["email"])
	assert.Equal(t, "081222222222", data["phone"])

	t.Logf("[PASS] PUT /api/v1/users/me - Success")
}

func Test_UpdateMyProfile_InvalidEmail(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"email": "invalid-email-format",
	}

	w := performRealRequest("PUT", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/me - Invalid Email")
}

func Test_UpdateMyProfile_UsernameTooShort(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "ab", // Kurang dari 3 karakter
	}

	w := performRealRequest("PUT", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/me - Username Too Short")
}

func Test_UpdateMyProfile_PhoneTooShort(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"phone": "0812345", // Kurang dari 10 karakter
	}

	w := performRealRequest("PUT", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/me - Phone Too Short")
}

func Test_UpdateMyProfile_DuplicateUsername(t *testing.T) {
	cleanupRealTestDB()

	// Buat user existing
	createRealUser("existinguser", "existing@example.com", "081111111111", "password123", models.RolePatient, true)

	// Buat user yang akan diupdate
	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "existinguser", // Username sudah ada
	}

	w := performRealRequest("PUT", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/me - Duplicate Username")
}

func Test_UpdateMyProfile_DuplicateEmail(t *testing.T) {
	cleanupRealTestDB()

	createRealUser("existinguser", "existing@example.com", "081111111111", "password123", models.RolePatient, true)
	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"email": "existing@example.com",
	}

	w := performRealRequest("PUT", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/me - Duplicate Email")
}

func Test_UpdateMyProfile_DuplicatePhone(t *testing.T) {
	cleanupRealTestDB()

	createRealUser("existinguser", "existing@example.com", "081111111111", "password123", models.RolePatient, true)
	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"phone": "081111111111",
	}

	w := performRealRequest("PUT", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/me - Duplicate Phone")
}

// ==================== PATCH /api/v1/users/me/change-password Tests ====================

func Test_ChangeMyPassword_Success(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"old_password": "password123",
		"new_password": "newpassword123",
	}

	w := performRealRequest("PATCH", "/api/v1/users/me/change-password", token, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/me/change-password - Success")
}

func Test_ChangeMyPassword_WrongOldPassword(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"old_password": "wrongpassword",
		"new_password": "newpassword123",
	}

	w := performRealRequest("PATCH", "/api/v1/users/me/change-password", token, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/me/change-password - Wrong Old Password")
}

func Test_ChangeMyPassword_NewPasswordTooShort(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"old_password": "password123",
		"new_password": "short", // Kurang dari 8 karakter
	}

	w := performRealRequest("PATCH", "/api/v1/users/me/change-password", token, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/me/change-password - New Password Too Short")
}

func Test_ChangeMyPassword_MissingOldPassword(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"new_password": "newpassword123",
	}

	w := performRealRequest("PATCH", "/api/v1/users/me/change-password", token, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/me/change-password - Missing Old Password")
}

// ==================== DELETE /api/v1/users/me Tests ====================

func Test_DeleteMyAccount_Success(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"password": "password123",
		"reason":   "No longer need this account",
	}

	w := performRealRequest("DELETE", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/me - Success")
}

func Test_DeleteMyAccount_WrongPassword(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"password": "wrongpassword",
	}

	w := performRealRequest("DELETE", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/me - Wrong Password")
}

func Test_DeleteMyAccount_MissingPassword(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"reason": "Taking a break",
	}

	w := performRealRequest("DELETE", "/api/v1/users/me", token, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/me - Missing Password")
}

// ==================== PATCH /api/v1/users/me/deactivate Tests ====================

func Test_DeactivateMyAccount_Success(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"password": "password123",
		"reason":   "Taking a break",
	}

	w := performRealRequest("PATCH", "/api/v1/users/me/deactivate", token, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/me/deactivate - Success")
}

func Test_DeactivateMyAccount_WrongPassword(t *testing.T) {
	cleanupRealTestDB()

	_, token, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"password": "wrongpassword",
	}

	w := performRealRequest("PATCH", "/api/v1/users/me/deactivate", token, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/me/deactivate - Wrong Password")
}

// ==================== POST /api/v1/users (Admin) Tests ====================

func Test_CreateUser_Success(t *testing.T) {
	cleanupRealTestDB()

	// Login sebagai admin
	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	// Request body seperti di Postman
	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "newuser", data["username"])
	assert.Equal(t, "newuser@example.com", data["email"])
	assert.Equal(t, "081234567890", data["phone"])
	assert.Equal(t, "patient", data["role"])

	t.Logf("[PASS] POST /api/v1/users - Success")
}

func Test_CreateUser_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	// Login sebagai patient
	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Forbidden for Patient")
}

func Test_CreateUser_InvalidEmail(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "invalid-email",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Invalid Email")
}

func Test_CreateUser_UsernameTooShort(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "ab",
		"email":    "newuser@example.com",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Username Too Short")
}

func Test_CreateUser_UsernameTooLong(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "thisusernameiswaytoolongandexceedsthemaximumlengthallowedof50characters",
		"email":    "newuser@example.com",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Username Too Long")
}

func Test_CreateUser_PasswordTooShort(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"phone":    "081234567890",
		"password": "short",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Password Too Short")
}

func Test_CreateUser_PhoneTooShort(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"phone":    "0812345",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Phone Too Short")
}

func Test_CreateUser_PhoneTooLong(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"phone":    "08123456789012345",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Phone Too Long")
}

func Test_CreateUser_InvalidRole(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "invalid_role",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Invalid Role")
}

func Test_CreateUser_DuplicateUsername(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	createRealUser("existinguser", "existing@example.com", "081222222222", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "existinguser",
		"email":    "newuser@example.com",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Duplicate Username")
}

func Test_CreateUser_DuplicateEmail(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	createRealUser("existinguser", "existing@example.com", "081222222222", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "existing@example.com",
		"phone":    "081234567890",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Duplicate Email")
}

func Test_CreateUser_DuplicatePhone(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	createRealUser("existinguser", "existing@example.com", "081222222222", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"phone":    "081222222222",
		"password": "password123",
		"role":     "patient",
	}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Duplicate Phone")
}

func Test_CreateUser_MissingRequiredFields(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{}

	w := performRealRequest("POST", "/api/v1/users", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/users - Missing Required Fields")
}

// ==================== GET /api/v1/users (List Users) Tests ====================

func Test_ListUsers_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	createRealUser("user1", "user1@example.com", "081222222222", "password123", models.RolePatient, true)
	createRealUser("user2", "user2@example.com", "081333333333", "password123", models.RoleDoctor, true)

	w := performRealRequest("GET", "/api/v1/users", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Debug: print response
	t.Logf("Response: %+v", response)

	// Response structure: { success, message, data: { data: [], meta: {} } }
	dataVal := response["data"]
	if dataVal == nil {
		t.Fatal("Response data is nil")
	}

	// data should be a map with "data" and "meta" fields
	dataMap, ok := dataVal.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data to be object, got: %T - value: %+v", dataVal, dataVal)
	}

	// Debug: print dataMap
	t.Logf("dataMap: %+v", dataMap)

	// Get the actual users array from data.data
	usersVal, ok := dataMap["data"]
	if !ok || usersVal == nil {
		t.Fatalf("Response data.data is nil or missing. dataMap keys: %+v", dataMap)
	}

	users := usersVal.([]interface{})
	assert.GreaterOrEqual(t, len(users), 3)
	t.Logf("[PASS] GET /api/v1/users - Success (Total: %d users)", len(users))
}

func Test_ListUsers_WithPagination(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	createRealUser("user1", "user1@example.com", "081222222222", "password123", models.RolePatient, true)
	createRealUser("user2", "user2@example.com", "081333333333", "password123", models.RoleDoctor, true)

	w := performRealRequest("GET", "/api/v1/users?page=1&page_size=2", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Response structure: { success, message, data: { data: [], meta: {} } }
	dataVal := response["data"]
	if dataVal == nil {
		t.Fatal("Response data is nil")
	}

	data, ok := dataVal.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data to be object, got: %T", dataVal)
	}

	meta, ok := data["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Response data.meta is nil or invalid")
	}

	assert.Equal(t, float64(1), meta["page"])
	assert.Equal(t, float64(2), meta["page_size"])

	t.Logf("[PASS] GET /api/v1/users - With Pagination")
}

func Test_ListUsers_WithSearch(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	createRealUser("searchuser", "searchuser@example.com", "081222222222", "password123", models.RolePatient, true)
	createRealUser("otheruser", "other@example.com", "081333333333", "password123", models.RoleDoctor, true)

	w := performRealRequest("GET", "/api/v1/users?search=searchuser", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/users - With Search")
}

func Test_ListUsers_WithRoleFilter(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	createRealUser("patient1", "patient1@example.com", "081222222222", "password123", models.RolePatient, true)
	createRealUser("doctor1", "doctor1@example.com", "081333333333", "password123", models.RoleDoctor, true)

	w := performRealRequest("GET", "/api/v1/users?role=patient", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/users - With Role Filter")
}

func Test_ListUsers_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)

	w := performRealRequest("GET", "/api/v1/users", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/users - Forbidden for Patient")
}

// ==================== GET /api/v1/users/:id Tests ====================

func Test_GetUserByID_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("GET", fmt.Sprintf("/api/v1/users/%d", user.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, user.Username, data["username"])
	assert.Equal(t, user.Email, data["email"])

	t.Logf("[PASS] GET /api/v1/users/:id - Success")
}

func Test_GetUserByID_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performRealRequest("GET", "/api/v1/users/99999", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/users/:id - Not Found")
}

func Test_GetUserByID_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("GET", fmt.Sprintf("/api/v1/users/%d", user.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/users/:id - Forbidden for Patient")
}

// ==================== PUT /api/v1/users/:id Tests ====================

func Test_UpdateUser_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "updateduser",
		"email":    "updated@example.com",
		"phone":    "081999999999",
	}

	w := performRealRequest("PUT", fmt.Sprintf("/api/v1/users/%d", user.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "updateduser", data["username"])
	assert.Equal(t, "updated@example.com", data["email"])

	t.Logf("[PASS] PUT /api/v1/users/:id - Success")
}

func Test_UpdateUser_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"username": "updateduser",
	}

	w := performRealRequest("PUT", "/api/v1/users/99999", adminToken, requestBody)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/:id - Not Found")
}

func Test_UpdateUser_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"username": "updateduser",
	}

	w := performRealRequest("PUT", fmt.Sprintf("/api/v1/users/%d", user.ID), patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PUT /api/v1/users/:id - Forbidden for Patient")
}

// ==================== DELETE /api/v1/users/:id (Soft Delete) Tests ====================

func Test_SoftDeleteUser_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", user.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/:id - Success (Soft Delete)")
}

func Test_SoftDeleteUser_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performRealRequest("DELETE", "/api/v1/users/99999", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/:id - Not Found")
}

func Test_SoftDeleteUser_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", user.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/:id - Forbidden for Patient")
}

// ==================== GET /api/v1/users/deleted Tests ====================

func Test_ListDeletedUsers_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("deleteduser", "deleted@example.com", "081234567890", "password123", models.RolePatient, true)

	// Soft delete user
	realTestDB.Delete(user)

	w := performRealRequest("GET", "/api/v1/users/deleted", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/users/deleted - Success")
}

func Test_ListDeletedUsers_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)

	w := performRealRequest("GET", "/api/v1/users/deleted", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/users/deleted - Forbidden for Patient")
}

// ==================== PATCH /api/v1/users/:id/restore Tests ====================

func Test_RestoreUser_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("deleteduser", "deleted@example.com", "081234567890", "password123", models.RolePatient, true)

	// Soft delete user first
	realTestDB.Delete(user)

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/restore", user.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/restore - Success")
}

func Test_RestoreUser_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performRealRequest("PATCH", "/api/v1/users/99999/restore", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/restore - Not Found")
}

func Test_RestoreUser_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("deleteduser", "deleted@example.com", "081234567890", "password123", models.RolePatient, true)

	realTestDB.Delete(user)

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/restore", user.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/restore - Forbidden for Patient")
}

// ==================== PATCH /api/v1/users/:id/reset-password Tests ====================

func Test_ResetPassword_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"new_password": "newpassword123",
	}

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/reset-password", user.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/reset-password - Success")
}

func Test_ResetPassword_PasswordTooShort(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"new_password": "short",
	}

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/reset-password", user.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/reset-password - Password Too Short")
}

func Test_ResetPassword_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"new_password": "newpassword123",
	}

	w := performRealRequest("PATCH", "/api/v1/users/99999/reset-password", adminToken, requestBody)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/reset-password - Not Found")
}

func Test_ResetPassword_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"new_password": "newpassword123",
	}

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/reset-password", user.ID), patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/reset-password - Forbidden for Patient")
}

// ==================== PATCH /api/v1/users/:id/activate Tests ====================

func Test_ActivateUser_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, false)

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/activate", user.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/activate - Success")
}

func Test_ActivateUser_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performRealRequest("PATCH", "/api/v1/users/99999/activate", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/activate - Not Found")
}

func Test_ActivateUser_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, false)

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/activate", user.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/activate - Forbidden for Patient")
}

// ==================== PATCH /api/v1/users/:id/deactivate Tests ====================

func Test_DeactivateUser_Success(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/deactivate", user.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/deactivate - Success")
}

func Test_DeactivateUser_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performRealRequest("PATCH", "/api/v1/users/99999/deactivate", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/deactivate - Not Found")
}

func Test_DeactivateUser_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("PATCH", fmt.Sprintf("/api/v1/users/%d/deactivate", user.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/users/:id/deactivate - Forbidden for Patient")
}

// ==================== DELETE /api/v1/users/:id/hard-delete Tests ====================

func Test_HardDeleteUser_Success_SuperAdmin(t *testing.T) {
	cleanupRealTestDB()

	_, superAdminToken, _ := createRealUser("superadmin", "superadmin@example.com", "081111111111", "password123", models.RoleSuperAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	// Soft delete first
	realTestDB.Delete(user)

	w := performRealRequest("DELETE", fmt.Sprintf("/api/v1/users/%d/hard-delete", user.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/:id/hard-delete - Success (Super Admin)")
}

func Test_HardDeleteUser_Forbidden_Admin(t *testing.T) {
	cleanupRealTestDB()

	_, adminToken, _ := createRealUser("admin", "admin@example.com", "081111111111", "password123", models.RoleAdmin, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("DELETE", fmt.Sprintf("/api/v1/users/%d/hard-delete", user.ID), adminToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/:id/hard-delete - Forbidden for Admin")
}

func Test_HardDeleteUser_NotFound(t *testing.T) {
	cleanupRealTestDB()

	_, superAdminToken, _ := createRealUser("superadmin", "superadmin@example.com", "081111111111", "password123", models.RoleSuperAdmin, true)

	w := performRealRequest("DELETE", "/api/v1/users/99999/hard-delete", superAdminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/:id/hard-delete - Not Found")
}

func Test_HardDeleteUser_Forbidden_Patient(t *testing.T) {
	cleanupRealTestDB()

	_, patientToken, _ := createRealUser("patient", "patient@example.com", "081111111111", "password123", models.RolePatient, true)
	user, _, _ := createRealUser("testuser", "test@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRealRequest("DELETE", fmt.Sprintf("/api/v1/users/%d/hard-delete", user.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/users/:id/hard-delete - Forbidden for Patient")
}
