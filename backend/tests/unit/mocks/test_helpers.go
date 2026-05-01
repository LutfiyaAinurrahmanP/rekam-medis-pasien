package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

// ============= User Model Builders =============

// NewTestUser creates a test user with default values
func NewTestUser() *models.User {
	return &models.User{
		ID:        1,
		Username:  "test_user",
		Email:     "test@example.com",
		Phone:     "08123456789",
		Password:  "hashed_password",
		Role:      models.RolePatient,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTestUserWithData creates a test user with custom data
func NewTestUserWithData(id uint, username, email, phone, role string, isActive bool) *models.User {
	return &models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Phone:     phone,
		Password:  "hashed_password",
		Role:      role,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ============= DTO Builders =============

// NewTestUserResponse creates a test user response with default values
func NewTestUserResponse() *dto.UserResponse {
	return &dto.UserResponse{
		ID:        1,
		Username:  "test_user",
		Email:     "test@example.com",
		Phone:     "08123456789",
		Role:      models.RolePatient,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTestUserResponseWithData creates a test user response with custom data
func NewTestUserResponseWithData(id uint, username, email, phone, role string, isActive bool) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        id,
		Username:  username,
		Email:     email,
		Phone:     phone,
		Role:      role,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewRegisterRequest creates a test register request
func NewRegisterRequest(username, email, phone, password string) *dto.RegisterRequest {
	return &dto.RegisterRequest{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: password,
		Role:     models.RolePatient,
	}
}

// NewCreateUserRequest creates a test create user request
func NewCreateUserRequest(username, email, phone, password, role string, isActive bool) *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: password,
		Role:     role,
		IsActive: &isActive,
	}
}

// NewLoginRequest creates a test login request
func NewLoginRequest(usernameOrEmail, password string) *dto.LoginRequest {
	return &dto.LoginRequest{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	}
}

// NewUpdateUserRequest creates a test update user request
func NewUpdateUserRequest() *dto.UpdateUserRequest {
	username := "updated_user"
	email := "updated@example.com"
	return &dto.UpdateUserRequest{
		Username: &username,
		Email:    &email,
	}
}

// NewChangePasswordRequest creates a test change password request
func NewChangePasswordRequest(oldPass, newPass string) *dto.ChangePasswordRequest {
	return &dto.ChangePasswordRequest{
		OldPassword: oldPass,
		NewPassword: newPass,
	}
}

// NewUserPaginationQuery creates a test pagination query
func NewUserPaginationQuery(page, pageSize int) *dto.UserPaginationQuery {
	return &dto.UserPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

// ============= Pointer Helpers =============

// PtrString returns a pointer to a string
func PtrString(s string) *string {
	return &s
}

// PtrBool returns a pointer to a bool
func PtrBool(b bool) *bool {
	return &b
}

// PtrUint returns a pointer to an uint
func PtrUint(u uint) *uint {
	return &u
}

// ============= Test List Builders =============

// NewTestUserList creates a list of test users
func NewTestUserList(count int) []models.User {
	users := make([]models.User, count)
	for i := 0; i < count; i++ {
		users[i] = *NewTestUserWithData(
			uint(i+1),
			"user_"+string(rune(i+1)),
			"user"+string(rune(i+1))+"@example.com",
			"0812345678"+string(rune(i+1)),
			models.RolePatient,
			true,
		)
	}
	return users
}

// NewTestUserResponseList creates a list of test user responses
func NewTestUserResponseList(count int) []dto.UserResponse {
	responses := make([]dto.UserResponse, count)
	for i := 0; i < count; i++ {
		responses[i] = *NewTestUserResponseWithData(
			uint(i+1),
			"user_"+string(rune(i+1)),
			"user"+string(rune(i+1))+"@example.com",
			"0812345678"+string(rune(i+1)),
			models.RolePatient,
			true,
		)
	}
	return responses
}

