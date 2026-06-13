package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockVitalSignService struct {
	mock.Mock
}

func (m *MockVitalSignService) List(query *dto.VitalSignPaginationQuery) (*dto.VitalSignListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VitalSignListResponse), args.Error(1)
}

func (m *MockVitalSignService) DeletedList(query *dto.VitalSignPaginationQuery) (*dto.VitalSignDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VitalSignDeletedListResponse), args.Error(1)
}

func (m *MockVitalSignService) FindByID(id uint) (*dto.VitalSignResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VitalSignResponse), args.Error(1)
}

func (m *MockVitalSignService) FindByIDUnscoped(id uint) (*dto.VitalSignResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VitalSignResponse), args.Error(1)
}

func (m *MockVitalSignService) Create(req *dto.CreateVitalSignRequest) (*dto.VitalSignResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VitalSignResponse), args.Error(1)
}

func (m *MockVitalSignService) Update(id uint, req *dto.UpdateVitalSignRequest) (*dto.VitalSignResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VitalSignResponse), args.Error(1)
}

func (m *MockVitalSignService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVitalSignService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVitalSignService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
