package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoomHandler_ListRooms_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := &dto.RoomListResponse{
		Data: []dto.RoomResponse{},
		Meta: dto.RoomPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("ListRooms", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListRooms(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_GetAvailableRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := &dto.RoomListResponse{
		Data: []dto.RoomResponse{},
		Meta: dto.RoomPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("GetAvailableRooms", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms/available?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAvailableRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_GetByOccupiedRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := &dto.RoomListResponse{
		Data: []dto.RoomResponse{},
		Meta: dto.RoomPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("GetOccupiedRooms", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms/occupied?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetByOccupiedRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_GetByActiveRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := &dto.RoomListResponse{
		Data: []dto.RoomResponse{},
		Meta: dto.RoomPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("GetActiveRooms", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms/active?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetByActiveRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_GetByInactiveRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := &dto.RoomListResponse{
		Data: []dto.RoomResponse{},
		Meta: dto.RoomPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("GetInactiveRooms", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms/inactive?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetByInactiveRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_DeletedListRooms_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := &dto.RoomDeletedListResponse{
		Data: []dto.DeletedRoomResponse{},
		Meta: dto.RoomPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("DeleteListRooms", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms/deleted?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.DeletedListRooms(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_GetRoomByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := mocks.NewTestRoomResponseWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockService.On("GetRoomByID", uint(1)).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetRoomByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_CreateRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	createReq := mocks.NewCreateRoomRequest("R001", 1, 1, 4, 4, 500000, true)
	body, _ := json.Marshal(createReq)

	expectedResponse := mocks.NewTestRoomResponseWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockService.On("CreateRoom", mock.AnythingOfType("*dto.CreateRoomRequest")).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateRoom(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRoomHandler_CreateRoom_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	createReq := mocks.NewCreateRoomRequest("R001", 1, 1, 4, 4, 500000, true)
	body, _ := json.Marshal(createReq)

	mockService.On("CreateRoom", mock.AnythingOfType("*dto.CreateRoomRequest")).Return(nil, errors.New("room number already exists"))

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateRoom(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRoomHandler_UpdateRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	updateReq := mocks.NewUpdateRoomRequest("R002", 2, 2, 6, 600000, true)
	body, _ := json.Marshal(updateReq)

	expectedResponse := mocks.NewTestRoomResponseWithData(1, "R002", 2, 2, 6, 6, 600000, true)

	mockService.On("UpdateRoom", uint(1), mock.AnythingOfType("*dto.UpdateRoomRequest")).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/rooms/1", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UpdateRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_ActivateRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := mocks.NewTestRoomResponseWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockService.On("ActivateRoom", uint(1)).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/rooms/1/activate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.ActivateRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_DeactivateRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	expectedResponse := mocks.NewTestRoomResponseWithData(1, "R001", 1, 1, 4, 4, 500000, false)

	mockService.On("DeactivateRoom", uint(1)).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/rooms/1/deactivate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeactivateRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_OccupyRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	occupyReq := dto.OccupyRoomRequest{Beds: 2}
	body, _ := json.Marshal(occupyReq)

	expectedResponse := mocks.NewTestRoomResponseWithData(1, "R001", 1, 1, 4, 2, 500000, true)

	mockService.On("OccupyRoom", uint(1), 2).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/rooms/1/occupy", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.OccupyRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_ReleaseRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	releaseReq := dto.ReleaseRoomRequest{Beds: 2}
	body, _ := json.Marshal(releaseReq)

	expectedResponse := mocks.NewTestRoomResponseWithData(1, "R001", 1, 1, 4, 4, 500000, true)

	mockService.On("ReleaseRoom", uint(1), 2).Return(expectedResponse, nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/rooms/1/release", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.ReleaseRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_SoftDeleteRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	mockService.On("SoftDeleteRoom", uint(1)).Return(nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/rooms/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDeleteRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_RestoreRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	mockService.On("RestoreRoom", uint(1)).Return(nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/rooms/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.RestoreRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomHandler_HardDeleteRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomService)

	mockService.On("HardDeleteRoom", uint(1)).Return(nil)

	h := handler.NewRoomHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/rooms/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDeleteRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
