package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockHospitalizationRepository struct {
	mock.Mock
}

func (m *MockHospitalizationRepository) List(query *dto.HospitalizationPaginationQuery) ([]models.Hospitalization, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Hospitalization), args.Get(1).(int64), args.Error(2)
}

func (m *MockHospitalizationRepository) DeletedList(query *dto.HospitalizationPaginationQuery) ([]models.Hospitalization, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Hospitalization), args.Get(1).(int64), args.Error(2)
}

func (m *MockHospitalizationRepository) FindByID(id uint) (*models.Hospitalization, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Hospitalization), args.Error(1)
}

func (m *MockHospitalizationRepository) FindByIDUnscoped(id uint) (*models.Hospitalization, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Hospitalization), args.Error(1)
}

func (m *MockHospitalizationRepository) Create(hospitalization *models.Hospitalization) error {
	args := m.Called(hospitalization)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) Update(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) Discharge(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) Transfer(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) Activate(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) Deactivate(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHospitalizationRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
