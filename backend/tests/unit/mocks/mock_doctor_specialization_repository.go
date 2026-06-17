package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockDoctorSpecializationRepository struct {
	mock.Mock
}

func (m *MockDoctorSpecializationRepository) List(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.DoctorSpecialization), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorSpecializationRepository) DeletedList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.DoctorSpecialization), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorSpecializationRepository) ActiveList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.DoctorSpecialization), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorSpecializationRepository) InactiveList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.DoctorSpecialization), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorSpecializationRepository) FindByID(id uint) (*models.DoctorSpecialization, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DoctorSpecialization), args.Error(1)
}

func (m *MockDoctorSpecializationRepository) Create(doctorSpecialization *models.DoctorSpecialization) error {
	args := m.Called(doctorSpecialization)
	return args.Error(0)
}

func (m *MockDoctorSpecializationRepository) Update(doctorSpecialization *models.DoctorSpecialization) error {
	args := m.Called(doctorSpecialization)
	return args.Error(0)
}

func (m *MockDoctorSpecializationRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationRepository) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationRepository) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	args := m.Called(name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDoctorSpecializationRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}
