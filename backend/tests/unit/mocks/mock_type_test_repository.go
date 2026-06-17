package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTypeTestRepository struct {
	mock.Mock
}

func (m *MockTypeTestRepository) List(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTest), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestRepository) DeletedList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTest), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestRepository) ActiveList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTest), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestRepository) InactiveList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.TypeTest), args.Get(1).(int64), args.Error(2)
}

func (m *MockTypeTestRepository) FindByID(id uint) (*models.TypeTest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TypeTest), args.Error(1)
}

func (m *MockTypeTestRepository) Create(typeTest *models.TypeTest) error {
	args := m.Called(typeTest)
	return args.Error(0)
}

func (m *MockTypeTestRepository) Update(typeTest *models.TypeTest) error {
	args := m.Called(typeTest)
	return args.Error(0)
}

func (m *MockTypeTestRepository) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestRepository) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	args := m.Called(name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTypeTestRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}
