package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMedicalRecordRepository struct {
	mock.Mock
}

func (m *MockMedicalRecordRepository) List(query *dto.MedicalRecordPaginationQuery) ([]models.MedicalRecord, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.MedicalRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicalRecordRepository) DeletedList(query *dto.MedicalRecordPaginationQuery) ([]models.MedicalRecord, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.MedicalRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicalRecordRepository) FindByID(id uint) (*models.MedicalRecord, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MedicalRecord), args.Error(1)
}

func (m *MockMedicalRecordRepository) FindByIDUnscoped(id uint) (*models.MedicalRecord, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MedicalRecord), args.Error(1)
}

func (m *MockMedicalRecordRepository) Create(record *models.MedicalRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockMedicalRecordRepository) Update(record *models.MedicalRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockMedicalRecordRepository) Finalize(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicalRecordRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicalRecordRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicalRecordRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
