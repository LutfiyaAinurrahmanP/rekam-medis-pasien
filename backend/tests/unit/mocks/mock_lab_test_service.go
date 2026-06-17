package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockLabTestService struct {
	mock.Mock
}

func (m *MockLabTestService) List(query *dto.LabTestPaginationQuery) (*dto.LabTestListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestListResponse), args.Error(1)
}

func (m *MockLabTestService) DeletedList(query *dto.LabTestPaginationQuery) (*dto.LabTestDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestDeletedListResponse), args.Error(1)
}

func (m *MockLabTestService) FindByID(id uint) (*dto.LabTestResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) FindByIDUnscoped(id uint) (*dto.LabTestResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) Create(req *dto.CreateLabTestRequest) (*dto.LabTestResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) Update(id uint, req *dto.UpdateLabTestRequest) (*dto.LabTestResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) CollectSample(id uint) (*dto.LabTestResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) Start(id uint) (*dto.LabTestResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) Complete(id uint, req *dto.CompleteLabTestRequest) (*dto.LabTestResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) Cancel(id uint) (*dto.LabTestResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LabTestResponse), args.Error(1)
}

func (m *MockLabTestService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLabTestService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLabTestService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
