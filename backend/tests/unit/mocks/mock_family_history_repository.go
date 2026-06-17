package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockFamilyHistoryRepository struct {
	mock.Mock
}

func (m *MockFamilyHistoryRepository) List(query *dto.FamilyHistoryPaginationQuery) ([]models.FamilyHistory, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.FamilyHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockFamilyHistoryRepository) FindByID(id uint) (*models.FamilyHistory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FamilyHistory), args.Error(1)
}

func (m *MockFamilyHistoryRepository) Create(history *models.FamilyHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockFamilyHistoryRepository) Update(history *models.FamilyHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockFamilyHistoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
