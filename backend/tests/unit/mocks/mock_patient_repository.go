package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/stretchr/testify/mock"
)

// MockPatientRepository is a mock implementation of PatientRepository interface
type MockPatientRepository struct {
	mock.Mock
}

func (m *MockPatientRepository) List(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Patient), args.Get(1).(int64), args.Error(2)
}

func (m *MockPatientRepository) DeleteList(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Patient), args.Get(1).(int64), args.Error(2)
}

func (m *MockPatientRepository) FindById(id uint) (*models.Patient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindByUserID(userID uint) (*models.Patient, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindByCode(code string) (*models.Patient, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) Create(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) Update(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}

// compile-time interface check
var _ repository.PatientRepository = (*MockPatientRepository)(nil)
