package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockFamilyHistoryService struct {
	mock.Mock
}

func (m *MockFamilyHistoryService) List(query *dto.FamilyHistoryPaginationQuery) (*dto.FamilyHistoryListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FamilyHistoryListResponse), args.Error(1)
}

func (m *MockFamilyHistoryService) FindByID(id uint) (*dto.FamilyHistoryResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FamilyHistoryResponse), args.Error(1)
}

func (m *MockFamilyHistoryService) Create(req *dto.CreateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FamilyHistoryResponse), args.Error(1)
}

func (m *MockFamilyHistoryService) Update(id uint, req *dto.UpdateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FamilyHistoryResponse), args.Error(1)
}

func (m *MockFamilyHistoryService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
