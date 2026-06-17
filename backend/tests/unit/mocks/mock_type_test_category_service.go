package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockTypeTestCategoryService struct {
	mock.Mock
}

func (m *MockTypeTestCategoryService) List(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestCategoryListResponse), args.Error(1)
}

func (m *MockTypeTestCategoryService) DeletedList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestCategoryDeletedListResponse), args.Error(1)
}

func (m *MockTypeTestCategoryService) ActiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestCategoryListResponse), args.Error(1)
}

func (m *MockTypeTestCategoryService) InactiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestCategoryListResponse), args.Error(1)
}

func (m *MockTypeTestCategoryService) FindByID(id uint) (*dto.TypeTestCategoryResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestCategoryResponse), args.Error(1)
}

func (m *MockTypeTestCategoryService) Create(req *dto.CreateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestCategoryResponse), args.Error(1)
}

func (m *MockTypeTestCategoryService) Update(id uint, req *dto.UpdateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TypeTestCategoryResponse), args.Error(1)
}

func (m *MockTypeTestCategoryService) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryService) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTypeTestCategoryService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
