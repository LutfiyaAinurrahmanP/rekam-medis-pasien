package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockSurgicalHistoryRepository struct {
	mock.Mock
}

func (m *MockSurgicalHistoryRepository) List(query *dto.SurgicalHistoryPaginationQuery) ([]models.SurgicalHistory, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.SurgicalHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockSurgicalHistoryRepository) FindByID(id uint) (*models.SurgicalHistory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SurgicalHistory), args.Error(1)
}

func (m *MockSurgicalHistoryRepository) Create(history *models.SurgicalHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockSurgicalHistoryRepository) Update(history *models.SurgicalHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockSurgicalHistoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
