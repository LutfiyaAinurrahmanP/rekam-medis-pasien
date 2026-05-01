package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService interface
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

func (m *MockUserService) ListUsers(query *dto.UserPaginationQuery) (*dto.UserListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserListResponse), args.Error(1)
}

func (m *MockUserService) DeleteListUsers(query *dto.UserPaginationQuery) (*dto.UserDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserDeletedListResponse), args.Error(1)
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
