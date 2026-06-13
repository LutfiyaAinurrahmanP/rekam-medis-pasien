package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	roomtypeservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room-type"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRoomTypeService() (*mocks.MockRoomTypeRepository, roomtypeservice.RoomTypeService) {
	mockRepo := new(mocks.MockRoomTypeRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := roomtypeservice.NewRoomTypeService(mockRepo, cfg)
	return mockRepo, service
}

func TestRoomTypeService_List_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("List", query).Return(roomTypes, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("DeletedList", query).Return(roomTypes, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	roomType := mocks.NewTestRoomTypeWithData(1, "ICU", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("FindByID", uint(1)).Return(roomType, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, roomType.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("room type not found"))

	res, err := service.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "room type not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	req := mocks.NewCreateRoomTypeRequest("ICU", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("IsNameExists", "ICU", mock.Anything).Return(false, nil)
	mockRepo.On("IsCodeExists", "ICU-01", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.RoomType")).Return(nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_Create_NameExists(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	req := mocks.NewCreateRoomTypeRequest("ICU", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("IsNameExists", "ICU", mock.Anything).Return(true, nil)

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "name already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_Create_CodeExists(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	req := mocks.NewCreateRoomTypeRequest("ICU", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("IsNameExists", "ICU", mock.Anything).Return(false, nil)
	mockRepo.On("IsCodeExists", "ICU-01", mock.Anything).Return(true, nil)

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "code already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	roomType := mocks.NewTestRoomTypeWithData(1, "ICU", "ICU-01", "Intensive Care Unit", true)
	req := mocks.NewUpdateRoomTypeRequest("ICU Updated", "ICU-02", "Updated", true)

	mockRepo.On("FindByID", uint(1)).Return(roomType, nil)
	mockRepo.On("IsNameExists", "ICU Updated", []uint{1}).Return(false, nil)
	mockRepo.On("IsCodeExists", "ICU-02", []uint{1}).Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.RoomType")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "ICU Updated", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	roomType := mocks.NewTestRoomTypeWithData(1, "ICU", "ICU-01", "Intensive Care Unit", true)

	mockRepo.On("FindByID", uint(1)).Return(roomType, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_ActiveList_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("ActiveList", query).Return(roomTypes, int64(2), nil)

	res, err := service.ActiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_InactiveList_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()
	query := mocks.NewRoomTypePaginationQuery(1, 10)
	roomTypes := mocks.NewTestRoomTypeList(2)

	mockRepo.On("InactiveList", query).Return(roomTypes, int64(2), nil)

	res, err := service.InactiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_Activate_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := service.Activate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoomTypeService_Deactivate_Success(t *testing.T) {
	mockRepo, service := setupTestRoomTypeService()

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := service.Deactivate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
