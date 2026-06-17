package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockRoomTypeService struct {
	mock.Mock
}

func (m *MockRoomTypeService) List(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomTypeListResponse), args.Error(1)
}

func (m *MockRoomTypeService) DeletedList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomTypeDeletedListResponse), args.Error(1)
}

func (m *MockRoomTypeService) FindByID(id uint) (*dto.RoomTypeResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomTypeResponse), args.Error(1)
}

func (m *MockRoomTypeService) Create(req *dto.CreateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomTypeResponse), args.Error(1)
}

func (m *MockRoomTypeService) Update(id uint, req *dto.UpdateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomTypeResponse), args.Error(1)
}

func (m *MockRoomTypeService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeService) ActiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomTypeListResponse), args.Error(1)
}

func (m *MockRoomTypeService) InactiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomTypeListResponse), args.Error(1)
}

func (m *MockRoomTypeService) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeService) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
