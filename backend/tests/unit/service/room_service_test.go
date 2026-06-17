package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	roomservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRoomService() (*mocks.MockRoomRepository, roomservice.RoomService) {
	mockRepo := new(mocks.MockRoomRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := roomservice.NewRoomService(mockRepo, cfg)
	return mockRepo, service
}

func TestRoomService_ListRooms_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("List", query).Return(rooms, int64(2), nil)

	res, err := service.ListRooms(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_DeleteListRooms_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	query := mocks.NewRoomPaginationQuery(1, 10)
	query.SortBy = "deleted_at"
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("DeleteList", query).Return(rooms, int64(2), nil)

	res, err := service.DeleteListRooms(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_GetAvailableRooms_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	query := mocks.NewRoomPaginationQuery(1, 10)
	query.SortBy = "available_beds"
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindAvailableRooms", query).Return(rooms, int64(2), nil)

	res, err := service.GetAvailableRooms(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_GetOccupiedRooms_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	query := mocks.NewRoomPaginationQuery(1, 10)
	query.SortBy = "available_beds"
	query.SortDir = "asc"
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindOccupiedRooms", query).Return(rooms, int64(2), nil)

	res, err := service.GetOccupiedRooms(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_GetActiveRooms_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindActiveRooms", query).Return(rooms, int64(2), nil)

	res, err := service.GetActiveRooms(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_GetInactiveRooms_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindInactiveRooms", query).Return(rooms, int64(2), nil)

	res, err := service.GetInactiveRooms(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_GetRoomByID_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)

	res, err := service.GetRoomByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, room.RoomNumber, res.RoomNumber)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_GetRoomByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestRoomService()

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("room not found"))

	res, err := service.GetRoomByID(999)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "room not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomService_CreateRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	req := mocks.NewCreateRoomRequest("R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("IsRoomNumberExists", "R001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Room")).Return(nil)

	res, err := service.CreateRoom(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.RoomNumber, res.RoomNumber)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_CreateRoom_RoomNumberExists(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	req := mocks.NewCreateRoomRequest("R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("IsRoomNumberExists", "R001", mock.Anything).Return(true, nil)

	res, err := service.CreateRoom(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "room number already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomService_UpdateRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, true)
	req := mocks.NewUpdateRoomRequest("R002", 2, 2, 6, 600000, false)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)
	mockRepo.On("IsRoomNumberExists", "R002", []uint{1}).Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Room")).Return(nil)

	res, err := service.UpdateRoom(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "R002", res.RoomNumber)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_ActivateRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, false)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Room")).Return(nil)

	res, err := service.ActivateRoom(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_DeactivateRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Room")).Return(nil)

	res, err := service.DeactivateRoom(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_OccupyRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Room")).Return(nil)

	res, err := service.OccupyRoom(1, 2)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, res.AvailableBeds)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_OccupyRoom_NotEnoughBeds(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 2, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)

	res, err := service.OccupyRoom(1, 3)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "not enough available beds", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomService_ReleaseRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 2, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Room")).Return(nil)

	res, err := service.ReleaseRoom(1, 2)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 4, res.AvailableBeds)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_ReleaseRoom_TooManyBeds(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 2, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)

	res, err := service.ReleaseRoom(1, 3)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot release more beds than capacity", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomService_SoftDeleteRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDeleteRoom(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_RestoreRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.RestoreRoom(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomService_HardDeleteRoom_Success(t *testing.T) {
	mockRepo, service := setupTestRoomService()

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDeleteRoom(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
