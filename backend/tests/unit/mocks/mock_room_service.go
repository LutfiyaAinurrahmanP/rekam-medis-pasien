package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockRoomService struct {
	mock.Mock
}

func (m *MockRoomService) ListRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomListResponse), args.Error(1)
}

func (m *MockRoomService) DeleteListRooms(query *dto.RoomPaginationQuery) (*dto.RoomDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomDeletedListResponse), args.Error(1)
}

func (m *MockRoomService) GetAvailableRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomListResponse), args.Error(1)
}

func (m *MockRoomService) GetOccupiedRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomListResponse), args.Error(1)
}

func (m *MockRoomService) GetActiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomListResponse), args.Error(1)
}

func (m *MockRoomService) GetInactiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomListResponse), args.Error(1)
}

func (m *MockRoomService) GetRoomByID(id uint) (*dto.RoomResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomResponse), args.Error(1)
}

func (m *MockRoomService) CreateRoom(req *dto.CreateRoomRequest) (*dto.RoomResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomResponse), args.Error(1)
}

func (m *MockRoomService) UpdateRoom(id uint, req *dto.UpdateRoomRequest) (*dto.RoomResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomResponse), args.Error(1)
}

func (m *MockRoomService) ActivateRoom(id uint) (*dto.RoomResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomResponse), args.Error(1)
}

func (m *MockRoomService) DeactivateRoom(id uint) (*dto.RoomResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomResponse), args.Error(1)
}

func (m *MockRoomService) OccupyRoom(id uint, beds int) (*dto.RoomResponse, error) {
	args := m.Called(id, beds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomResponse), args.Error(1)
}

func (m *MockRoomService) ReleaseRoom(id uint, beds int) (*dto.RoomResponse, error) {
	args := m.Called(id, beds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoomResponse), args.Error(1)
}

func (m *MockRoomService) SoftDeleteRoom(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomService) RestoreRoom(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomService) HardDeleteRoom(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
