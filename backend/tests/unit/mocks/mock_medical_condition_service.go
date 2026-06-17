package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockMedicalConditionService struct {
	mock.Mock
}

func (m *MockMedicalConditionService) List(query *dto.MedicalConditionPaginationQuery) (*dto.MedicalConditionListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalConditionListResponse), args.Error(1)
}

func (m *MockMedicalConditionService) FindByID(id uint) (*dto.MedicalConditionResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalConditionResponse), args.Error(1)
}

func (m *MockMedicalConditionService) Create(req *dto.CreateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalConditionResponse), args.Error(1)
}

func (m *MockMedicalConditionService) Update(id uint, req *dto.UpdateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalConditionResponse), args.Error(1)
}

func (m *MockMedicalConditionService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
