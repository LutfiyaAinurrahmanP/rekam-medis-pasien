package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockDepartmentRepository is a mock implementation of DepartmentRepository interface
type MockDepartmentRepository struct {
	mock.Mock
}

func (m *MockDepartmentRepository) List(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Department), args.Get(1).(int64), args.Error(2)
}

func (m *MockDepartmentRepository) DeleteList(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Department), args.Get(1).(int64), args.Error(2)
}

func (m *MockDepartmentRepository) FindById(id uint) (*models.Department, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepository) FindByName(name string) (*models.Department, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepository) FindByCode(code string) (*models.Department, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepository) Create(department *models.Department) error {
	args := m.Called(department)
	return args.Error(0)
}

func (m *MockDepartmentRepository) Update(department *models.Department) error {
	args := m.Called(department)
	return args.Error(0)
}

func (m *MockDepartmentRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDepartmentRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDepartmentRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDepartmentRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}
