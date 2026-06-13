package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockDoctorRepository struct {
	mock.Mock
}

func (m *MockDoctorRepository) List(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Doctor), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorRepository) DeleteList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Doctor), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorRepository) ActiveList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Doctor), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorRepository) InactiveList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Doctor), args.Get(1).(int64), args.Error(2)
}

func (m *MockDoctorRepository) FindByID(id uint) (*models.Doctor, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Doctor), args.Error(1)
}

func (m *MockDoctorRepository) FindByUserID(userID uint) (*models.Doctor, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Doctor), args.Error(1)
}

func (m *MockDoctorRepository) FindByDepartmentID(departmentID uint) (*models.Doctor, error) {
	args := m.Called(departmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Doctor), args.Error(1)
}

func (m *MockDoctorRepository) Create(doctor *models.Doctor) error {
	args := m.Called(doctor)
	return args.Error(0)
}

func (m *MockDoctorRepository) Update(doctor *models.Doctor) error {
	args := m.Called(doctor)
	return args.Error(0)
}

func (m *MockDoctorRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorRepository) IsEmployeeIDExists(employeeID string, excludeID ...uint) (bool, error) {
	args := m.Called(employeeID, excludeID)
	return args.Bool(0), args.Error(1)
}
