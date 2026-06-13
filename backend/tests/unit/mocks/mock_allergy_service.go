package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockAllergyService struct {
	mock.Mock
}

func (m *MockAllergyService) List(query *dto.AllergyPaginationQuery) (*dto.AllergyListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AllergyListResponse), args.Error(1)
}

func (m *MockAllergyService) FindByID(id uint) (*dto.AllergyResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AllergyResponse), args.Error(1)
}

func (m *MockAllergyService) Create(req *dto.CreateAllergyRequest) (*dto.AllergyResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AllergyResponse), args.Error(1)
}

func (m *MockAllergyService) Update(id uint, req *dto.UpdateAllergyRequest) (*dto.AllergyResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AllergyResponse), args.Error(1)
}

func (m *MockAllergyService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
