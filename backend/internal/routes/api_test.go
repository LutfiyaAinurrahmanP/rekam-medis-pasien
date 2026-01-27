package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockUserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) ListUsers(query *dto.PaginationQuery) (*dto.UserListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserListResponse), args.Error(1)
}

func (m *MockUserService) DeleteListUsers(query *dto.PaginationQuery) (*dto.DeletedUserListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DeletedUserListResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) SoftDeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) HardDeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) RestoreUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) ChangePassword(id uint, req *dto.ChangePasswordRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockUserService) ResetPassword(id uint, newPassword string) error {
	args := m.Called(id, newPassword)
	return args.Error(0)
}

func (m *MockUserService) ActivateUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) DeactivateUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) VerifyPasswordForDeletion(id uint, password string) error {
	args := m.Called(id, password)
	return args.Error(0)
}

// Test helpers
type TestResult struct {
	Name     string
	Passed   bool
	Duration time.Duration
	Error    string
}

type TestSuite struct {
	Name    string
	Results []TestResult
	Start   time.Time
}

// Global test statistics
type GlobalTestStats struct {
	sync.Mutex
	TotalSuites   int
	TotalTests    int
	PassedTests   int
	FailedTests   int
	TotalDuration time.Duration
}

var globalStats = &GlobalTestStats{}

func (ts *TestSuite) AddResult(name string, passed bool, duration time.Duration, err string) {
	ts.Results = append(ts.Results, TestResult{
		Name:     name,
		Passed:   passed,
		Duration: duration,
		Error:    err,
	})
}

func (ts *TestSuite) PrintSummary() {
	elapsed := time.Since(ts.Start)
	passed := 0
	failed := 0

	fmt.Printf("\n %s\n", ts.Name)

	for _, result := range ts.Results {
		if result.Passed {
			passed++
			fmt.Printf("   ✓ %s (%v)\n", result.Name, result.Duration)
		} else {
			failed++
			fmt.Printf("   ✗ %s (%v)\n", result.Name, result.Duration)
			if result.Error != "" {
				fmt.Printf("     Error: %s\n", result.Error)
			}
		}
	}

	// Store stats for global summary (only count once per suite)
	globalStats.Lock()
	if len(ts.Results) > 0 {
		globalStats.TotalSuites++
		globalStats.TotalTests += passed + failed
		globalStats.PassedTests += passed
		globalStats.FailedTests += failed
		globalStats.TotalDuration += elapsed
	}
	globalStats.Unlock()
}

func PrintGlobalSummary() {
	globalStats.Lock()
	defer globalStats.Unlock()

	// Use ASCII dashes instead of Unicode for Windows compatibility
	fmt.Println("\n" + strings.Repeat("-", 65))

	passedSuites := globalStats.TotalSuites
	if globalStats.FailedTests > 0 {
		passedSuites = globalStats.TotalSuites - 1 // Estimate failed suites
	}

	// Format like Jest
	fmt.Printf("Test Suites: ")
	if globalStats.FailedTests > 0 {
		fmt.Printf("%d failed, ", globalStats.TotalSuites-passedSuites)
	}
	fmt.Printf("%d passed, %d total\n", passedSuites, globalStats.TotalSuites)

	fmt.Printf("Tests:       ")
	if globalStats.FailedTests > 0 {
		fmt.Printf("%d failed, ", globalStats.FailedTests)
	}
	fmt.Printf("%d passed, %d total\n", globalStats.PassedTests, globalStats.TotalTests)

	fmt.Printf("Snapshots:   0 total\n")
	fmt.Printf("Time:        %.3f s\n", globalStats.TotalDuration.Seconds())
	fmt.Println("Ran all test suites.")
	fmt.Println()
}

func setupTestRouter() (*gin.Engine, *MockUserService, *config.Config) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key",
			ExpiredTime: 24 * 60 * time.Minute, // 24 hours
		},
	}

	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	routeConfig := &RouteConfig{
		UserHandler: userHandler,
		Config:      cfg,
	}

	api := router.Group("/api/v1")
	setupAPIRouter(api, routeConfig)

	return router, mockService, cfg
}

func generateTestToken(userID uint, username, email, role string, cfg *config.Config) string {
	token, _, _ := utils.GenerateToken(userID, username, email, role, cfg.JWT.Secret, cfg.JWT.ExpiredTime)
	return token
}

func performRequest(router *gin.Engine, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ============================================================================
// Test Cases for /users/me endpoints
// ============================================================================

func TestGetMyProfile(t *testing.T) {
	suite := &TestSuite{Name: "GET /api/v1/users/me - Get My Profile", Start: time.Now()}

	t.Run("Success - Get my profile", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		expectedResponse := &dto.UserResponse{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Phone:    "08123456789",
			Role:     models.RolePatient,
			IsActive: true,
		}

		mockService.On("GetUserByID", uint(1)).Return(expectedResponse, nil)

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "GET", "/api/v1/users/me", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should return 200 and user profile", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Missing authorization token", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		w := performRequest(router, "GET", "/api/v1/users/me", "", nil)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 when token missing", passed, time.Since(start), "")
	})

	t.Run("Error - Invalid token format", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", "InvalidToken")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 for invalid token format", passed, time.Since(start), "")
	})

	t.Run("Error - User not found", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("GetUserByID", uint(1)).Return(nil, fmt.Errorf("user not found"))

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "GET", "/api/v1/users/me", token, nil)

		passed := assert.Equal(t, http.StatusNotFound, w.Code)
		suite.AddResult("should return 404 when user not found", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	suite.PrintSummary()
}

func TestUpdateMyProfile(t *testing.T) {
	suite := &TestSuite{Name: "PUT /api/v1/users/me - Update My Profile", Start: time.Now()}

	t.Run("Success - Update profile", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		newEmail := "newemail@example.com"
		updateReq := dto.UpdateUserRequest{
			Email: &newEmail,
		}

		expectedResponse := &dto.UserResponse{
			ID:       1,
			Username: "testuser",
			Email:    newEmail,
			Phone:    "08123456789",
			Role:     models.RolePatient,
			IsActive: true,
		}

		mockService.On("UpdateUser", uint(1), mock.AnythingOfType("*dto.UpdateUserRequest")).Return(expectedResponse, nil)

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PUT", "/api/v1/users/me", token, updateReq)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully update profile", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Try to update role (forbidden)", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		role := models.RoleAdmin
		updateReq := dto.UpdateUserRequest{
			Role: &role,
		}

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PUT", "/api/v1/users/me", token, updateReq)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject role update attempt", passed, time.Since(start), "")
	})

	t.Run("Error - Try to update is_active (forbidden)", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		isActive := false
		updateReq := dto.UpdateUserRequest{
			IsActive: &isActive,
		}

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PUT", "/api/v1/users/me", token, updateReq)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject is_active update attempt", passed, time.Since(start), "")
	})

	t.Run("Error - Missing authorization", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		updateReq := dto.UpdateUserRequest{}
		w := performRequest(router, "PUT", "/api/v1/users/me", "", updateReq)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 when not authenticated", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestChangeMyPassword(t *testing.T) {
	suite := &TestSuite{Name: "PATCH /api/v1/users/me/change-password - Change My Password", Start: time.Now()}

	t.Run("Success - Change password", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		changePassReq := dto.ChangePasswordRequest{
			OldPassword: "oldpassword123",
			NewPassword: "newpassword123",
		}

		mockService.On("ChangePassword", uint(1), mock.AnythingOfType("*dto.ChangePasswordRequest")).Return(nil)

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/me/change-password", token, changePassReq)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully change password", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Wrong old password", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		changePassReq := dto.ChangePasswordRequest{
			OldPassword: "wrongpassword",
			NewPassword: "newpassword123",
		}

		mockService.On("ChangePassword", uint(1), mock.AnythingOfType("*dto.ChangePasswordRequest")).Return(fmt.Errorf("incorrect old password"))

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/me/change-password", token, changePassReq)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should reject wrong old password", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Missing authorization", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		changePassReq := dto.ChangePasswordRequest{
			OldPassword: "oldpassword123",
			NewPassword: "newpassword123",
		}

		w := performRequest(router, "PATCH", "/api/v1/users/me/change-password", "", changePassReq)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 when not authenticated", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestDeleteMyAccount(t *testing.T) {
	suite := &TestSuite{Name: "DELETE /api/v1/users/me - Delete My Account", Start: time.Now()}

	t.Run("Success - Delete account", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		deleteReq := dto.DeleteAccountRequest{
			Password: "correctpassword",
		}

		mockService.On("VerifyPasswordForDeletion", uint(1), "correctpassword").Return(nil)
		mockService.On("SoftDeleteUser", uint(1)).Return(nil)

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/me", token, deleteReq)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully delete account", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Wrong password", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		deleteReq := dto.DeleteAccountRequest{
			Password: "wrongpassword",
		}

		mockService.On("VerifyPasswordForDeletion", uint(1), "wrongpassword").Return(fmt.Errorf("incorrect password"))

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/me", token, deleteReq)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should reject wrong password", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Missing authorization", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		deleteReq := dto.DeleteAccountRequest{
			Password: "password123",
		}

		w := performRequest(router, "DELETE", "/api/v1/users/me", "", deleteReq)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 when not authenticated", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestDeactivateMyAccount(t *testing.T) {
	suite := &TestSuite{Name: "PATCH /api/v1/users/me/deactivate - Deactivate My Account", Start: time.Now()}

	t.Run("Success - Deactivate account", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		deactivateReq := dto.DeactivateAccountRequest{
			Password: "correctpassword",
		}

		mockService.On("VerifyPasswordForDeletion", uint(1), "correctpassword").Return(nil)
		mockService.On("DeactivateUser", uint(1)).Return(nil)

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/me/deactivate", token, deactivateReq)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully deactivate account", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Wrong password", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		deactivateReq := dto.DeactivateAccountRequest{
			Password: "wrongpassword",
		}

		mockService.On("VerifyPasswordForDeletion", uint(1), "wrongpassword").Return(fmt.Errorf("incorrect password"))

		token := generateTestToken(1, "testuser", "test@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/me/deactivate", token, deactivateReq)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should reject wrong password", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Missing authorization", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		deactivateReq := dto.DeactivateAccountRequest{
			Password: "password123",
		}

		w := performRequest(router, "PATCH", "/api/v1/users/me/deactivate", "", deactivateReq)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 when not authenticated", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

// ============================================================================
// Test Cases for Admin endpoints
// ============================================================================

func TestCreateUser(t *testing.T) {
	suite := &TestSuite{Name: "POST /api/v1/users - Create User (Admin)", Start: time.Now()}

	t.Run("Success - Admin creates user", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		createReq := dto.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Phone:    "08199999999",
			Password: "password123",
			Role:     models.RolePatient,
		}

		expectedResponse := &dto.UserResponse{
			ID:       2,
			Username: "newuser",
			Email:    "newuser@example.com",
			Phone:    "08199999999",
			Role:     models.RolePatient,
			IsActive: true,
		}

		mockService.On("CreateUser", mock.AnythingOfType("*dto.CreateUserRequest")).Return(expectedResponse, nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "POST", "/api/v1/users", token, createReq)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully create user as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to create user", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		createReq := dto.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Phone:    "08199999999",
			Password: "password123",
			Role:     models.RolePatient,
		}

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "POST", "/api/v1/users", token, createReq)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - Duplicate username", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		createReq := dto.CreateUserRequest{
			Username: "existinguser",
			Email:    "new@example.com",
			Phone:    "08199999999",
			Password: "password123",
			Role:     models.RolePatient,
		}

		mockService.On("CreateUser", mock.AnythingOfType("*dto.CreateUserRequest")).Return(nil, fmt.Errorf("username already exists"))

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "POST", "/api/v1/users", token, createReq)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should reject duplicate username", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Missing authorization", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		createReq := dto.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Phone:    "08199999999",
			Password: "password123",
			Role:     models.RolePatient,
		}

		w := performRequest(router, "POST", "/api/v1/users", "", createReq)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 when not authenticated", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestListUsers(t *testing.T) {
	suite := &TestSuite{Name: "GET /api/v1/users - List Users (Admin)", Start: time.Now()}

	t.Run("Success - Admin lists users", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		expectedResponse := &dto.UserListResponse{
			Data: []dto.UserResponse{
				{ID: 1, Username: "user1", Email: "user1@example.com"},
				{ID: 2, Username: "user2", Email: "user2@example.com"},
			},
			Meta: dto.PaginationMeta{
				Page:       1,
				PageSize:   10,
				TotalItems: 2,
				TotalPages: 1,
			},
		}

		mockService.On("ListUsers", mock.AnythingOfType("*dto.PaginationQuery")).Return(expectedResponse, nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "GET", "/api/v1/users?page=1&page_size=10", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully list users as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to list users", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "GET", "/api/v1/users", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - Missing authorization", func(t *testing.T) {
		start := time.Now()
		router, _, _ := setupTestRouter()

		w := performRequest(router, "GET", "/api/v1/users", "", nil)

		passed := assert.Equal(t, http.StatusUnauthorized, w.Code)
		suite.AddResult("should return 401 when not authenticated", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestGetDeletedUsers(t *testing.T) {
	suite := &TestSuite{Name: "GET /api/v1/users/deleted - Get Deleted Users (Admin)", Start: time.Now()}

	t.Run("Success - Admin gets deleted users", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		expectedResponse := &dto.DeletedUserListResponse{
			Data: []dto.DeletedUserResponse{
				{ID: 1, Username: "deleted1", Email: "deleted1@example.com"},
			},
			Meta: dto.PaginationMeta{
				Page:       1,
				PageSize:   10,
				TotalItems: 1,
				TotalPages: 1,
			},
		}

		mockService.On("DeleteListUsers", mock.AnythingOfType("*dto.PaginationQuery")).Return(expectedResponse, nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "GET", "/api/v1/users/deleted", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully get deleted users as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to access", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "GET", "/api/v1/users/deleted", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestGetUserByID(t *testing.T) {
	suite := &TestSuite{Name: "GET /api/v1/users/:id - Get User By ID (Admin)", Start: time.Now()}

	t.Run("Success - Admin gets user by ID", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		expectedResponse := &dto.UserResponse{
			ID:       2,
			Username: "targetuser",
			Email:    "target@example.com",
			Role:     models.RolePatient,
		}

		mockService.On("GetUserByID", uint(2)).Return(expectedResponse, nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "GET", "/api/v1/users/2", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully get user by ID as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - User not found", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("GetUserByID", uint(999)).Return(nil, fmt.Errorf("user not found"))

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "GET", "/api/v1/users/999", token, nil)

		passed := assert.Equal(t, http.StatusNotFound, w.Code)
		suite.AddResult("should return 404 when user not found", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "GET", "/api/v1/users/invalid", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return 400 for invalid ID", passed, time.Since(start), "")
	})

	t.Run("Error - Non-admin tries to access", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "GET", "/api/v1/users/2", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestUpdateUser(t *testing.T) {
	suite := &TestSuite{Name: "PUT /api/v1/users/:id - Update User (Admin)", Start: time.Now()}

	t.Run("Success - Admin updates user", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		newEmail := "updated@example.com"
		updateReq := dto.UpdateUserRequest{
			Email: &newEmail,
		}

		expectedResponse := &dto.UserResponse{
			ID:       2,
			Username: "targetuser",
			Email:    newEmail,
			Role:     models.RolePatient,
		}

		mockService.On("UpdateUser", uint(2), mock.AnythingOfType("*dto.UpdateUserRequest")).Return(expectedResponse, nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PUT", "/api/v1/users/2", token, updateReq)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully update user as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to update", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		newEmail := "updated@example.com"
		updateReq := dto.UpdateUserRequest{
			Email: &newEmail,
		}

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PUT", "/api/v1/users/2", token, updateReq)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - User not found", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		newEmail := "updated@example.com"
		updateReq := dto.UpdateUserRequest{
			Email: &newEmail,
		}

		mockService.On("UpdateUser", uint(999), mock.AnythingOfType("*dto.UpdateUserRequest")).Return(nil, fmt.Errorf("user not found"))

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PUT", "/api/v1/users/999", token, updateReq)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return error when user not found", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	suite.PrintSummary()
}

func TestSoftDeleteUser(t *testing.T) {
	suite := &TestSuite{Name: "DELETE /api/v1/users/:id - Soft Delete User (Admin)", Start: time.Now()}

	t.Run("Success - Admin soft deletes user", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("SoftDeleteUser", uint(2)).Return(nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/2", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully soft delete user as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to delete", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/2", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/invalid", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return 400 for invalid ID", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestRestoreUser(t *testing.T) {
	suite := &TestSuite{Name: "PATCH /api/v1/users/:id/restore - Restore User (Admin)", Start: time.Now()}

	t.Run("Success - Admin restores user", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("RestoreUser", uint(2)).Return(nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/restore", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully restore user as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to restore", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/restore", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - User not found", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("RestoreUser", uint(999)).Return(fmt.Errorf("user not found"))

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/999/restore", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return error when user not found", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	suite.PrintSummary()
}

func TestResetPassword(t *testing.T) {
	suite := &TestSuite{Name: "PATCH /api/v1/users/:id/reset-password - Reset Password (Admin)", Start: time.Now()}

	t.Run("Success - Admin resets password", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		resetReq := dto.ResetPasswordRequest{
			NewPassword: "newpassword123",
		}

		mockService.On("ResetPassword", uint(2), "newpassword123").Return(nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/reset-password", token, resetReq)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully reset password as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to reset", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		resetReq := dto.ResetPasswordRequest{
			NewPassword: "newpassword123",
		}

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/reset-password", token, resetReq)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		resetReq := dto.ResetPasswordRequest{
			NewPassword: "newpassword123",
		}

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/invalid/reset-password", token, resetReq)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return 400 for invalid ID", passed, time.Since(start), "")
	})

	suite.PrintSummary()
}

func TestActivateUser(t *testing.T) {
	suite := &TestSuite{Name: "PATCH /api/v1/users/:id/activate - Activate User (Admin)", Start: time.Now()}

	t.Run("Success - Admin activates user", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("ActivateUser", uint(2)).Return(nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/activate", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully activate user as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to activate", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/activate", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - User not found", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("ActivateUser", uint(999)).Return(fmt.Errorf("user not found"))

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/999/activate", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return error when user not found", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	suite.PrintSummary()
}

func TestDeactivateUser(t *testing.T) {
	suite := &TestSuite{Name: "PATCH /api/v1/users/:id/deactivate - Deactivate User (Admin)", Start: time.Now()}

	t.Run("Success - Admin deactivates user", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		// Admin deactivates user with ID 2 (not themselves)
		mockService.On("DeactivateUser", uint(2)).Return(nil)

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/deactivate", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully deactivate user as admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Non-admin tries to deactivate", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/2/deactivate", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/invalid/deactivate", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return 400 for invalid ID", passed, time.Since(start), "")
	})

	t.Run("Error - User not found", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("DeactivateUser", uint(999)).Return(fmt.Errorf("user not found"))

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "PATCH", "/api/v1/users/999/deactivate", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return error when user not found", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	suite.PrintSummary()
}

// ============================================================================
// Test Cases for Super Admin endpoints
// ============================================================================

func TestHardDeleteUser(t *testing.T) {
	suite := &TestSuite{Name: "DELETE /api/v1/users/:id/hard-delete - Hard Delete User (SuperAdmin)", Start: time.Now()}

	t.Run("Success - SuperAdmin hard deletes user", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("HardDeleteUser", uint(2)).Return(nil)

		token := generateTestToken(1, "superadmin", "superadmin@example.com", models.RoleSuperAdmin, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/2/hard-delete", token, nil)

		passed := assert.Equal(t, http.StatusOK, w.Code)
		suite.AddResult("should successfully hard delete user as super admin", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Admin tries to hard delete (forbidden)", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "admin", "admin@example.com", models.RoleAdmin, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/2/hard-delete", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject admin user (requires super admin)", passed, time.Since(start), "")
	})

	t.Run("Error - Patient tries to hard delete", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "patient", "patient@example.com", models.RolePatient, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/2/hard-delete", token, nil)

		passed := assert.Equal(t, http.StatusForbidden, w.Code)
		suite.AddResult("should reject non-super-admin user", passed, time.Since(start), "")
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		start := time.Now()
		router, _, cfg := setupTestRouter()

		token := generateTestToken(1, "superadmin", "superadmin@example.com", models.RoleSuperAdmin, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/invalid/hard-delete", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return 400 for invalid ID", passed, time.Since(start), "")
	})

	t.Run("Error - User not found", func(t *testing.T) {
		start := time.Now()
		router, mockService, cfg := setupTestRouter()

		mockService.On("HardDeleteUser", uint(999)).Return(fmt.Errorf("user not found"))

		token := generateTestToken(1, "superadmin", "superadmin@example.com", models.RoleSuperAdmin, cfg)
		w := performRequest(router, "DELETE", "/api/v1/users/999/hard-delete", token, nil)

		passed := assert.Equal(t, http.StatusBadRequest, w.Code)
		suite.AddResult("should return error when user not found", passed, time.Since(start), "")

		mockService.AssertExpectations(t)
	})

	suite.PrintSummary()
}

// TestMain runs before all tests and after all tests
func TestMain(m *testing.M) {
	// Reset global stats before running tests
	globalStats = &GlobalTestStats{}

	// Run all tests
	code := m.Run()

	// Print global summary
	PrintGlobalSummary()

	// Exit with the test result code
	fmt.Printf("\x1b[0m") // Reset color
	if code == 0 {
		fmt.Println()
	}
}
