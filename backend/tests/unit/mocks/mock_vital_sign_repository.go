package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockVitalSignRepository struct {
	mock.Mock
}

func (m *MockVitalSignRepository) List(query *dto.VitalSignPaginationQuery) ([]models.VitalSign, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.VitalSign), args.Get(1).(int64), args.Error(2)
}

func (m *MockVitalSignRepository) DeletedList(query *dto.VitalSignPaginationQuery) ([]models.VitalSign, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.VitalSign), args.Get(1).(int64), args.Error(2)
}

func (m *MockVitalSignRepository) FindByID(id uint) (*models.VitalSign, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VitalSign), args.Error(1)
}

func (m *MockVitalSignRepository) FindByIDUnscoped(id uint) (*models.VitalSign, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VitalSign), args.Error(1)
}

func (m *MockVitalSignRepository) FindByMedicalRecordID(recordID uint) (*models.VitalSign, error) {
	args := m.Called(recordID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VitalSign), args.Error(1)
}

func (m *MockVitalSignRepository) Create(vitalSign *models.VitalSign) error {
	args := m.Called(vitalSign)
	return args.Error(0)
}

func (m *MockVitalSignRepository) Update(vitalSign *models.VitalSign) error {
	args := m.Called(vitalSign)
	return args.Error(0)
}

func (m *MockVitalSignRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVitalSignRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVitalSignRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
