package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============= Test Cases: GetMyProfile (Self-Owned) =============

func TestGetMyProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	expectedUser := &dto.UserResponse{
		ID:        1,
		Username:  "john_doe",
		Email:     "john@example.com",
		Phone:     "08123456789",
		Role:      "patient",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetUserByID", uint(1)).Return(expectedUser, nil)

	h := handler.NewUserHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.GetMyProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Profile retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestGetMyProfile_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)
	h := handler.NewUserHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetMyProfile(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["message"])
	assert.Equal(t, "user not authenticated", response["error"])
}

func TestGetMyProfile_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("GetUserByID", uint(999)).Return(nil, errors.New("user not found"))

	h := handler.NewUserHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(999))

	h.GetMyProfile(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User not found", response["message"])
	assert.Equal(t, "user not found", response["error"])
}

// ============= Test Cases: UpdateMyProfile =============

func TestUpdateMyProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	updatedUser := &dto.UserResponse{
		ID:       1,
		Username: "john_updated",
		Email:    "john_new@example.com",
		Role:     "patient",
		IsActive: true,
	}

	mockService.On("UpdateUser", uint(1), mock.Anything).Return(updatedUser, nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/users/me", strings.NewReader(`{
		"username": "john_updated",
		"email": "john_new@example.com"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.UpdateMyProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateMyProfile_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/users/me", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.UpdateMyProfile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateMyProfile_CannotChangeRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/users/me", strings.NewReader(`{
		"role": "admin"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.UpdateMyProfile(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ============= Test Cases: ChangeMyPassword =============

func TestChangeMyPassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("ChangePassword", uint(1), mock.Anything).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PATCH", "/users/me/change-password", strings.NewReader(`{
		"old_password": "oldpass123",
		"new_password": "newpass456"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.ChangeMyPassword(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestChangeMyPassword_InvalidOldPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("ChangePassword", uint(1), mock.Anything).
		Return(errors.New("old password is incorrect"))

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PATCH", "/users/me/change-password", strings.NewReader(`{
		"old_password": "wrongpass",
		"new_password": "newpass456"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.ChangeMyPassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: DeleteMyAccount =============

func TestDeleteMyAccount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("VerifyPasswordForDeletion", uint(1), mock.Anything).Return(nil)
	mockService.On("SoftDeleteUser", uint(1)).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("DELETE", "/users/me", strings.NewReader(`{
		"password": "password123"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.DeleteMyAccount(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: DeactivateMyAccount =============

func TestDeactivateMyAccount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("VerifyPasswordForDeletion", uint(1), mock.Anything).Return(nil)
	mockService.On("DeactivateUser", uint(1)).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PATCH", "/users/me/deactivate", strings.NewReader(`{
		"password": "password123"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.DeactivateMyAccount(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: CreateUser (Admin) =============

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	newUser := &dto.UserResponse{
		ID:       2,
		Username: "drsmith",
		Email:    "drsmith@example.com",
		Role:     "doctor",
		IsActive: true,
	}

	mockService.On("CreateUser", mock.Anything).Return(newUser, nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{
		"username": "drsmith",
		"email": "drsmith@example.com",
		"phone": "08123456791",
		"password": "password123",
		"role": "doctor"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("CreateUser", mock.Anything).
		Return(nil, errors.New("username already exists"))

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{
		"username": "existinguser",
		"email": "new@example.com",
		"phone": "08123456799",
		"password": "password123"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateUser(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: ListUsers (Admin) =============

func TestListUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	users := &dto.UserListResponse{
		Data: []dto.UserResponse{
			{ID: 1, Username: "user1", Role: "patient", IsActive: true},
			{ID: 2, Username: "user2", Role: "doctor", IsActive: true},
		},
		Meta: dto.UserPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	mockService.On("ListUsers", mock.Anything).Return(users, nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/users?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: DeleteListUsers (Admin) =============

func TestDeleteListUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	deletedUsers := &dto.UserDeletedListResponse{
		Data: []dto.DeletedUserResponse{
			{ID: 5, Username: "deleteduser", Email: "deleted@example.com", IsActive: false},
		},
		Meta: dto.UserPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 1,
			TotalPages: 1,
		},
	}

	mockService.On("DeleteListUsers", mock.Anything).Return(deletedUsers, nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/users/deleted?page=1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.DeleteListUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: GetUserByID (Admin) =============

func TestGetUserByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	user := &dto.UserResponse{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Role:     "patient",
		IsActive: true,
	}

	mockService.On("GetUserByID", uint(1)).Return(user, nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/users/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.GetUserByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("GetUserByID", uint(999)).Return(nil, errors.New("user not found"))

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/users/999", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	h.GetUserByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: UpdateUser (Admin) =============

func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	updated := &dto.UserResponse{
		ID:       1,
		Username: "johndoe_updated",
		Role:     "doctor",
		IsActive: false,
	}

	mockService.On("UpdateUser", uint(1), mock.Anything).Return(updated, nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/users/1", strings.NewReader(`{
		"username": "johndoe_updated",
		"role": "doctor"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: SoftDeleteUser (Admin) =============

func TestSoftDeleteUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("SoftDeleteUser", uint(1)).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("DELETE", "/users/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.SoftDeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: RestoreUser (Admin) =============

func TestRestoreUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("RestoreUser", uint(1)).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PATCH", "/users/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.RestoreUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: ResetPassword (Admin) =============

func TestResetPassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("ResetPassword", uint(1), mock.Anything).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PATCH", "/users/1/reset-password", strings.NewReader(`{
		"new_password": "newpass123"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.ResetPassword(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: ActivateUser (Admin) =============

func TestActivateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("ActivateUser", uint(1)).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PATCH", "/users/1/activate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.ActivateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: DeactivateUser (Admin) =============

func TestDeactivateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("DeactivateUser", uint(1)).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PATCH", "/users/1/deactivate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.DeactivateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: HardDeleteUser (Super Admin) =============

func TestHardDeleteUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)

	mockService.On("HardDeleteUser", uint(1)).Return(nil)

	h := handler.NewUserHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("DELETE", "/users/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	h.HardDeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
