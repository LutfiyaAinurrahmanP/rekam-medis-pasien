package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockTypeTestService struct {
	mock.Mock
}

func (m *MockTypeTestService) List(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestListResponse), args.Error(1)
}

func (m *MockTypeTestService) DeletedList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestDeletedListResponse), args.Error(1)
}

func (m *MockTypeTestService) ActiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestListResponse), args.Error(1)
}

func (m *MockTypeTestService) InactiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestListResponse), args.Error(1)
}

func (m *MockTypeTestService) FindByID(id uint) (*dto.TypeTestResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestResponse), args.Error(1)
}

func (m *MockTypeTestService) Create(req *dto.CreateTypeTestRequest) (*dto.TypeTestResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestResponse), args.Error(1)
}

func (m *MockTypeTestService) Update(id uint, req *dto.UpdateTypeTestRequest) (*dto.TypeTestResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestResponse), args.Error(1)
}

func (m *MockTypeTestService) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestService) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
