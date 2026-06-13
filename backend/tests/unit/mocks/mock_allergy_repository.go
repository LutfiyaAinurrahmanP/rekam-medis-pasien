package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAllergyRepository struct {
	mock.Mock
}

func (m *MockAllergyRepository) List(query *dto.AllergyPaginationQuery) ([]models.Allergy, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Allergy), args.Get(1).(int64), args.Error(2)
}

func (m *MockAllergyRepository) FindByID(id uint) (*models.Allergy, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Allergy), args.Error(1)
}

func (m *MockAllergyRepository) Create(allergy *models.Allergy) error {
	args := m.Called(allergy)
	return args.Error(0)
}

func (m *MockAllergyRepository) Update(allergy *models.Allergy) error {
	args := m.Called(allergy)
	return args.Error(0)
}

func (m *MockAllergyRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
