package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockDoctorSpecializationService struct {
	mock.Mock
}

func (m *MockDoctorSpecializationService) List(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorSpecializationListResponse), args.Error(1)
}

func (m *MockDoctorSpecializationService) DeletedList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorSpecializationDeletedListResponse), args.Error(1)
}

func (m *MockDoctorSpecializationService) ActiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorSpecializationListResponse), args.Error(1)
}

func (m *MockDoctorSpecializationService) InactiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorSpecializationListResponse), args.Error(1)
}

func (m *MockDoctorSpecializationService) FindByID(id uint) (*dto.DoctorSpecializationResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorSpecializationResponse), args.Error(1)
}

func (m *MockDoctorSpecializationService) Create(req *dto.CreateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorSpecializationResponse), args.Error(1)
}

func (m *MockDoctorSpecializationService) Update(id uint, req *dto.UpdateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorSpecializationResponse), args.Error(1)
}

func (m *MockDoctorSpecializationService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationService) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorSpecializationService) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
