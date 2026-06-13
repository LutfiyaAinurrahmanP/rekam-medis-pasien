package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockMedicalHistoryService struct {
	mock.Mock
}

func (m *MockMedicalHistoryService) List(query *dto.MedicalHistoryPaginationQuery) (*dto.MedicalHistoryListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalHistoryListResponse), args.Error(1)
}

func (m *MockMedicalHistoryService) FindByID(id uint) (*dto.MedicalHistoryDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalHistoryDetailResponse), args.Error(1)
}
