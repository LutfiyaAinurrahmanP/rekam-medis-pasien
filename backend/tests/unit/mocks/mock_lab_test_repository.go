package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockLabTestRepository struct {
	mock.Mock
}

func (m *MockLabTestRepository) List(query *dto.LabTestPaginationQuery) ([]models.LabTest, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.LabTest), args.Get(1).(int64), args.Error(2)
}

func (m *MockLabTestRepository) DeletedList(query *dto.LabTestPaginationQuery) ([]models.LabTest, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.LabTest), args.Get(1).(int64), args.Error(2)
}

func (m *MockLabTestRepository) FindByID(id uint) (*models.LabTest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LabTest), args.Error(1)
}

func (m *MockLabTestRepository) FindByIDUnscoped(id uint) (*models.LabTest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LabTest), args.Error(1)
}

func (m *MockLabTestRepository) Create(labTest *models.LabTest) error {
	args := m.Called(labTest)
	return args.Error(0)
}

func (m *MockLabTestRepository) Update(labTest *models.LabTest) error {
	args := m.Called(labTest)
	return args.Error(0)
}

func (m *MockLabTestRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLabTestRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLabTestRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
