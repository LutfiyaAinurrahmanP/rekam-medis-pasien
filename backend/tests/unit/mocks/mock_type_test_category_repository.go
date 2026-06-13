package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTypeTestCategoryRepository struct {
	mock.Mock
}

func (m *MockTypeTestCategoryRepository) List(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTestCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestCategoryRepository) DeletedList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTestCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestCategoryRepository) ActiveList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTestCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestCategoryRepository) InactiveList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTestCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestCategoryRepository) FindByID(id uint) (*models.TypeTestCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TypeTestCategory), args.Error(1)
}

func (m *MockTypeTestCategoryRepository) Create(category *models.TypeTestCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockTypeTestCategoryRepository) Update(category *models.TypeTestCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockTypeTestCategoryRepository) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryRepository) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	args := m.Called(name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTypeTestCategoryRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}
