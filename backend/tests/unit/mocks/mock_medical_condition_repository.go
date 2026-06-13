package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMedicalConditionRepository struct {
	mock.Mock
}

func (m *MockMedicalConditionRepository) List(query *dto.MedicalConditionPaginationQuery) ([]models.MedicalCondition, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.MedicalCondition), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicalConditionRepository) FindByID(id uint) (*models.MedicalCondition, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MedicalCondition), args.Error(1)
}

func (m *MockMedicalConditionRepository) Create(condition *models.MedicalCondition) error {
	args := m.Called(condition)
	return args.Error(0)
}

func (m *MockMedicalConditionRepository) Update(condition *models.MedicalCondition) error {
	args := m.Called(condition)
	return args.Error(0)
}

func (m *MockMedicalConditionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
