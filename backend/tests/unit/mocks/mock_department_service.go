package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

// MockDepartmentService is a mock implementation of DepartmentService interface
type MockDepartmentService struct {
	mock.Mock
}

func (m *MockDepartmentService) ListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentListResponse), args.Error(1)
}

func (m *MockDepartmentService) DeleteListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentDeletedListResponse), args.Error(1)
}

func (m *MockDepartmentService) GetDepartmentByID(id uint) (*dto.DepartmentResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentResponse), args.Error(1)
}

func (m *MockDepartmentService) CreateDepartment(req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentResponse), args.Error(1)
}

func (m *MockDepartmentService) UpdateDepartment(id uint, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DepartmentResponse), args.Error(1)
}

func (m *MockDepartmentService) SoftDeleteDepartment(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDepartmentService) RestoreDepartment(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDepartmentService) HardDeleteDepartment(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
