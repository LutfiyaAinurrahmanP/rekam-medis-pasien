package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockPrescriptionRepository struct {
	mock.Mock
}

func (m *MockPrescriptionRepository) List(query *dto.PrescriptionPaginationQuery) ([]models.Prescription, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Prescription), args.Get(1).(int64), args.Error(2)
}

func (m *MockPrescriptionRepository) DeletedList(query *dto.PrescriptionPaginationQuery) ([]models.Prescription, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Prescription), args.Get(1).(int64), args.Error(2)
}

func (m *MockPrescriptionRepository) FindByID(id uint) (*models.Prescription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Prescription), args.Error(1)
}

func (m *MockPrescriptionRepository) FindByIDUnscoped(id uint) (*models.Prescription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Prescription), args.Error(1)
}

func (m *MockPrescriptionRepository) Create(prescription *models.Prescription) error {
	args := m.Called(prescription)
	return args.Error(0)
}

func (m *MockPrescriptionRepository) Update(prescription *models.Prescription) error {
	args := m.Called(prescription)
	return args.Error(0)
}

func (m *MockPrescriptionRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPrescriptionRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPrescriptionRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPrescriptionRepository) FindItemByID(id uint) (*models.PrescriptionItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PrescriptionItem), args.Error(1)
}

func (m *MockPrescriptionRepository) CreateItem(item *models.PrescriptionItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockPrescriptionRepository) UpdateItem(item *models.PrescriptionItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockPrescriptionRepository) DeleteItem(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
