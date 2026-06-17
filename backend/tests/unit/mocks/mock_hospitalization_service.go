package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockHospitalizationService struct {
	mock.Mock
}

func (m *MockHospitalizationService) List(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationListResponse), args.Error(1)
}

func (m *MockHospitalizationService) DeletedList(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationDeletedListResponse), args.Error(1)
}

func (m *MockHospitalizationService) FindByID(id uint) (*dto.HospitalizationResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) FindByIDUnscoped(id uint) (*dto.HospitalizationResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) Create(req *dto.CreateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) Update(id uint, req *dto.UpdateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) Discharge(id uint, req *dto.DischargeHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) Transfer(id uint, req *dto.TransferHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) Activate(id uint) (*dto.HospitalizationResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) Deactivate(id uint) (*dto.HospitalizationResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.HospitalizationResponse), args.Error(1)
}

func (m *MockHospitalizationService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHospitalizationService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHospitalizationService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
