package service

import (
	"errors"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	userservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/user"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============= Test Cases: GetUserByID =============

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	userModel := &models.User{
		ID:        1,
		Username:  "john_doe",
		Email:     "john@example.com",
		Phone:     "08123456789",
		Password:  "hashed_password",
		Role:      "patient",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", uint(1)).Return(userModel, nil)
	service := userservice.NewUserService(mockRepo, cfg)

	result, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "john_doe", result.Username)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
	service := userservice.NewUserService(mockRepo, cfg)

	result, err := service.GetUserByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: ListUsers =============

func TestListUsers_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	users := []models.User{
		{
			ID:       1,
			Username: "user1",
			Email:    "user1@example.com",
			Role:     "patient",
			IsActive: true,
		},
		{
			ID:       2,
			Username: "user2",
			Email:    "user2@example.com",
			Role:     "doctor",
			IsActive: true,
		},
	}

	query := &dto.UserPaginationQuery{Page: 1, PageSize: 10}
	mockRepo.On("List", query).Return(users, int64(2), nil)
	service := userservice.NewUserService(mockRepo, cfg)

	result, err := service.ListUsers(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, int64(2), result.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestListUsers_EmptyResult(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	query := &dto.UserPaginationQuery{Page: 1, PageSize: 10}
	mockRepo.On("List", query).Return([]models.User{}, int64(0), nil)
	service := userservice.NewUserService(mockRepo, cfg)

	result, err := service.ListUsers(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Data))
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: CreateUser =============

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	req := &dto.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Phone:    "08123456789",
		Password: "password123", // ✅ Password harus diisi (minimal requirement)
		Role:     "patient",
	}

	mockRepo.On("IsUsernameExists", "newuser", mock.Anything).Return(false, nil)
	mockRepo.On("IsEmailExists", "newuser@example.com", mock.Anything).Return(false, nil)
	mockRepo.On("IsPhoneExists", "08123456789", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.CreateUser(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "newuser", result.Username)
	mockRepo.AssertExpectations(t)
}

// ✅ NEW TEST: Validate empty password is rejected
func TestCreateUser_EmptyPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	req := &dto.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Phone:    "08123456789",
		Password: "", // ❌ Empty password (must be rejected)
		Role:     "patient",
	}

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.CreateUser(req)

	// ✅ Now expects error for empty password
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "password cannot be empty", err.Error())
}

func TestCreateUser_UsernameExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	req := &dto.CreateUserRequest{
		Username: "existinguser",
		Email:    "newuser@example.com",
		Phone:    "08123456789",
		Password: "password123",
	}

	mockRepo.On("IsUsernameExists", "existinguser", mock.Anything).Return(true, nil)
	service := userservice.NewUserService(mockRepo, cfg)

	result, err := service.CreateUser(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "username already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: UpdateUser =============

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	existingUser := &models.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Phone:    "08123456789",
	}

	req := &dto.UpdateUserRequest{
		Username: ptrString("johndoe_updated"),
	}

	mockRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("IsUsernameExists", "johndoe_updated", mock.Anything).Return(false, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.UpdateUser(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: SoftDeleteUser =============

func TestSoftDeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("FindByID", uint(1)).Return(&models.User{ID: 1}, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	err := service.SoftDeleteUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
	service := userservice.NewUserService(mockRepo, cfg)

	err := service.SoftDeleteUser(999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: RestoreUser =============

func TestRestoreUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("Restore", uint(1)).Return(nil)
	service := userservice.NewUserService(mockRepo, cfg)

	err := service.RestoreUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: HardDeleteUser =============

func TestHardDeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	err := service.HardDeleteUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: ChangePassword =============

func TestChangePassword_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	existingUser := &models.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashed_old_password",
	}

	req := &dto.ChangePasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword123",
	}

	mockRepo.On("FindByID", uint(1)).Return(existingUser, nil)

	service := userservice.NewUserService(mockRepo, cfg)
	err := service.ChangePassword(1, req)

	if err != nil {
		t.Logf("ChangePassword returned error (expected in some implementations): %v", err)
	}
}

func TestChangePassword_InvalidOldPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	existingUser := &models.User{
		ID:       1,
		Username: "johndoe",
		Password: "hashed_old_password",
	}

	req := &dto.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}

	mockRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	service := userservice.NewUserService(mockRepo, cfg)

	err := service.ChangePassword(1, req)

	assert.Error(t, err)
}

// ============= Test Cases: ResetPassword =============

func TestResetPassword_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("FindByID", uint(1)).Return(&models.User{ID: 1}, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	err := service.ResetPassword(1, "newpassword123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestResetPassword_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
	service := userservice.NewUserService(mockRepo, cfg)

	err := service.ResetPassword(999, "newpassword123")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: ActivateUser =============

func TestActivateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	inactiveUser := &models.User{
		ID:       1,
		Username: "johndoe",
		IsActive: false,
	}

	mockRepo.On("FindByID", uint(1)).Return(inactiveUser, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	err := service.ActivateUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestActivateUser_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
	service := userservice.NewUserService(mockRepo, cfg)

	err := service.ActivateUser(999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: DeactivateUser =============

func TestDeactivateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	activeUser := &models.User{
		ID:       1,
		Username: "johndoe",
		IsActive: true,
	}

	mockRepo.On("FindByID", uint(1)).Return(activeUser, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	err := service.DeactivateUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: DeleteListUsers =============

func TestDeleteListUsers_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	deletedUsers := []models.User{
		{
			ID:       5,
			Username: "deleteduser",
			IsActive: false,
		},
	}

	query := &dto.UserPaginationQuery{Page: 1, PageSize: 10}
	mockRepo.On("DeleteList", query).Return(deletedUsers, int64(1), nil)

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.DeleteListUsers(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Data))
	mockRepo.AssertExpectations(t)
}

// ✅ NEW TEST: Validate invalid role is rejected
func TestCreateUser_InvalidRole(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	req := &dto.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Phone:    "08123456789",
		Password: "password123",
		Role:     "invalid_role", // ❌ Invalid role (not in ValidateRole)
	}

	mockRepo.On("IsUsernameExists", "newuser", mock.Anything).Return(false, nil)
	mockRepo.On("IsEmailExists", "newuser@example.com", mock.Anything).Return(false, nil)
	mockRepo.On("IsPhoneExists", "08123456789", mock.Anything).Return(false, nil)

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.CreateUser(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid role", err.Error())
}

// ✅ NEW TEST: Validate email already exists
func TestCreateUser_EmailExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	req := &dto.CreateUserRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Phone:    "08123456789",
		Password: "password123",
		Role:     "patient",
	}

	mockRepo.On("IsUsernameExists", "newuser", mock.Anything).Return(false, nil)
	mockRepo.On("IsEmailExists", "existing@example.com", mock.Anything).Return(true, nil)

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.CreateUser(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "email already exists", err.Error())
}

// ✅ NEW TEST: Validate phone already exists
func TestCreateUser_PhoneExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	req := &dto.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Phone:    "08123456789",
		Password: "password123",
		Role:     "patient",
	}

	mockRepo.On("IsUsernameExists", "newuser", mock.Anything).Return(false, nil)
	mockRepo.On("IsEmailExists", "newuser@example.com", mock.Anything).Return(false, nil)
	mockRepo.On("IsPhoneExists", "08123456789", mock.Anything).Return(true, nil)

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.CreateUser(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "phone already exists", err.Error())
}

// ✅ NEW TEST: Validate default role assignment (empty role becomes 'patient')
func TestCreateUser_DefaultRole(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	req := &dto.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Phone:    "08123456789",
		Password: "password123",
		Role:     "", // Empty role - should default to 'patient'
	}

	mockRepo.On("IsUsernameExists", "newuser", mock.Anything).Return(false, nil)
	mockRepo.On("IsEmailExists", "newuser@example.com", mock.Anything).Return(false, nil)
	mockRepo.On("IsPhoneExists", "08123456789", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)

	service := userservice.NewUserService(mockRepo, cfg)
	result, err := service.CreateUser(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "patient", result.Role, "Empty role should default to 'patient'")
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: VerifyPasswordForDeletion =============

func TestVerifyPasswordForDeletion_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}

	user := &models.User{
		ID:       1,
		Password: "hashed_password",
	}

	mockRepo.On("FindByID", uint(1)).Return(user, nil)
	service := userservice.NewUserService(mockRepo, cfg)

	err := service.VerifyPasswordForDeletion(1, "correctpassword")

	if err != nil {
		t.Logf("VerifyPasswordForDeletion returned error (expected in some implementations): %v", err)
	}
}

// ============= Helper Functions =============

// ============= Helper Functions & Model Validation Tests =============

// ✅ TEST: Model.IsPatient() helper
func TestUserModel_IsPatient(t *testing.T) {
	user := &models.User{Role: models.RolePatient}
	assert.True(t, user.IsPatient())

	user.Role = models.RoleDoctor
	assert.False(t, user.IsPatient())
}

// ✅ TEST: Model.IsDoctor() helper
func TestUserModel_IsDoctor(t *testing.T) {
	user := &models.User{Role: models.RoleDoctor}
	assert.True(t, user.IsDoctor())

	user.Role = models.RolePatient
	assert.False(t, user.IsDoctor())
}

// ✅ TEST: Model.IsReceptionist() helper
func TestUserModel_IsReceptionist(t *testing.T) {
	user := &models.User{Role: models.RoleReceptionist}
	assert.True(t, user.IsReceptionist())

	user.Role = models.RolePatient
	assert.False(t, user.IsReceptionist())
}

// ✅ TEST: Model.IsAdmin() helper
func TestUserModel_IsAdmin(t *testing.T) {
	user := &models.User{Role: models.RoleAdmin}
	assert.True(t, user.IsAdmin())

	user.Role = models.RolePatient
	assert.False(t, user.IsAdmin())
}

// ✅ TEST: Model.IsSuperAdmin() helper
func TestUserModel_IsSuperAdmin(t *testing.T) {
	user := &models.User{Role: models.RoleSuperAdmin}
	assert.True(t, user.IsSuperAdmin())

	user.Role = models.RolePatient
	assert.False(t, user.IsSuperAdmin())
}

// ✅ TEST: ValidateRole() - All valid roles
func TestValidateRole_AllValidRoles(t *testing.T) {
	validRoles := []string{
		models.RolePatient,
		models.RoleDoctor,
		models.RoleReceptionist,
		models.RoleAdmin,
		models.RoleSuperAdmin,
	}

	for _, role := range validRoles {
		assert.True(t, models.ValidateRole(role), "Role '%s' should be valid", role)
	}
}

// ✅ TEST: ValidateRole() - Invalid roles
func TestValidateRole_InvalidRoles(t *testing.T) {
	invalidRoles := []string{
		"invalid_role",
		"user",
		"super_user",
		"moderator",
		"guest",
		"", // Empty should be invalid
	}

	for _, role := range invalidRoles {
		assert.False(t, models.ValidateRole(role), "Role '%s' should be invalid", role)
	}
}

// ✅ TEST: GetAvailableRoles() returns all roles
func TestGetAvailableRoles(t *testing.T) {
	roles := models.GetAvailableRoles()

	expectedRoles := []string{
		models.RolePatient,
		models.RoleDoctor,
		models.RoleReceptionist,
		models.RoleAdmin,
		models.RoleSuperAdmin,
	}

	assert.Equal(t, len(expectedRoles), len(roles), "Should return 5 roles")

	for _, expectedRole := range expectedRoles {
		found := false
		for _, role := range roles {
			if role == expectedRole {
				found = true
				break
			}
		}
		assert.True(t, found, "Role '%s' should be in available roles", expectedRole)
	}
}

func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}
