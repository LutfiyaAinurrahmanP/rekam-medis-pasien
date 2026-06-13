package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRoomTypeRepository struct {
	mock.Mock
}

func (m *MockRoomTypeRepository) List(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.RoomType), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomTypeRepository) DeletedList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.RoomType), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomTypeRepository) FindByID(id uint) (*models.RoomType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomType), args.Error(1)
}

func (m *MockRoomTypeRepository) Create(roomType *models.RoomType) error {
	args := m.Called(roomType)
	return args.Error(0)
}

func (m *MockRoomTypeRepository) Update(roomType *models.RoomType) error {
	args := m.Called(roomType)
	return args.Error(0)
}

func (m *MockRoomTypeRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	args := m.Called(name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoomTypeRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoomTypeRepository) ActiveList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.RoomType), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomTypeRepository) InactiveList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.RoomType), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomTypeRepository) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoomTypeRepository) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
