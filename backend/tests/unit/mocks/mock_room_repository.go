package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) List(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomRepository) DeleteList(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomRepository) FindAvailableRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomRepository) FindOccupiedRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomRepository) FindActiveRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomRepository) FindInactiveRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomRepository) FindByID(id uint) (*models.Room, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomRepository) Create(room *models.Room) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *MockRoomRepository) Update(room *models.Room) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *MockRoomRepository) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomRepository) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomRepository) IsRoomNumberExists(roomNumber string, excludeID ...uint) (bool, error) {
	args := m.Called(roomNumber, excludeID)
	return args.Bool(0), args.Error(1)
}
