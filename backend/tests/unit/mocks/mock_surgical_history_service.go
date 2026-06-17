package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockSurgicalHistoryService struct {
	mock.Mock
}

func (m *MockSurgicalHistoryService) List(query *dto.SurgicalHistoryPaginationQuery) (*dto.SurgicalHistoryListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SurgicalHistoryListResponse), args.Error(1)
}

func (m *MockSurgicalHistoryService) FindByID(id uint) (*dto.SurgicalHistoryResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SurgicalHistoryResponse), args.Error(1)
}

func (m *MockSurgicalHistoryService) Create(req *dto.CreateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SurgicalHistoryResponse), args.Error(1)
}

func (m *MockSurgicalHistoryService) Update(id uint, req *dto.UpdateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SurgicalHistoryResponse), args.Error(1)
}

func (m *MockSurgicalHistoryService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
