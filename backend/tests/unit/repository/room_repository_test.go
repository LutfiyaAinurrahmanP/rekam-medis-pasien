package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoomRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("List", query).Return(rooms, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_DeleteList(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("DeleteList", query).Return(rooms, int64(2), nil)

	res, total, err := mockRepo.DeleteList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_FindAvailableRooms(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindAvailableRooms", query).Return(rooms, int64(2), nil)

	res, total, err := mockRepo.FindAvailableRooms(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_FindOccupiedRooms(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindOccupiedRooms", query).Return(rooms, int64(2), nil)

	res, total, err := mockRepo.FindOccupiedRooms(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_FindActiveRooms(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindActiveRooms", query).Return(rooms, int64(2), nil)

	res, total, err := mockRepo.FindActiveRooms(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_FindInactiveRooms(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	query := mocks.NewRoomPaginationQuery(1, 10)
	rooms := mocks.NewTestRoomList(2)

	mockRepo.On("FindInactiveRooms", query).Return(rooms, int64(2), nil)

	res, total, err := mockRepo.FindInactiveRooms(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("FindByID", uint(1)).Return(room, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, room.RoomNumber, res.RoomNumber)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("room not found"))

	res, err := mockRepo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "room not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	room := mocks.NewTestRoomWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockRepo.On("Create", room).Return(nil)

	err := mockRepo.Create(room)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)
	room := mocks.NewTestRoomWithData(1, "R001 Updated", 1, 1, 4, 4, 500000, true)

	mockRepo.On("Update", room).Return(nil)

	err := mockRepo.Update(room)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := mockRepo.Activate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := mockRepo.Deactivate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_IsRoomNumberExists_True(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("IsRoomNumberExists", "R001", mock.Anything).Return(true, nil)

	exists, err := mockRepo.IsRoomNumberExists("R001")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestRoomRepository_IsRoomNumberExists_False(t *testing.T) {
	mockRepo := new(mocks.MockRoomRepository)

	mockRepo.On("IsRoomNumberExists", "R001", mock.Anything).Return(false, nil)

	exists, err := mockRepo.IsRoomNumberExists("R001")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}
