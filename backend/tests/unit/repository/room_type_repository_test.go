package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoomTypeRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("List", query).Return(roomTypes, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("DeletedList", query).Return(roomTypes, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_ActiveList(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("ActiveList", query).Return(roomTypes, int64(2), nil)

	res, total, err := mockRepo.ActiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_InactiveList(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("InactiveList", query).Return(roomTypes, int64(2), nil)

	res, total, err := mockRepo.InactiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	roomType := mocks.NewTestRoomTypeWithData(1, "ICU", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("FindByID", uint(1)).Return(roomType, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, roomType.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("room type not found"))

	res, err := mockRepo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "room type not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	roomType := mocks.NewTestRoomTypeWithData(1, "ICU", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("Create", roomType).Return(nil)

	err := mockRepo.Create(roomType)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	roomType := mocks.NewTestRoomTypeWithData(1, "ICU Updated", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("Update", roomType).Return(nil)

	err := mockRepo.Update(roomType)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_IsNameExists_True(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("IsNameExists", "ICU", mock.Anything).Return(true, nil)

	exists, err := mockRepo.IsNameExists("ICU")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_IsCodeExists_False(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("IsCodeExists", "ICU-01", mock.Anything).Return(false, nil)

	exists, err := mockRepo.IsCodeExists("ICU-01")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := mockRepo.Activate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeRepository_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockRoomTypeRepository)

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := mockRepo.Deactivate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
