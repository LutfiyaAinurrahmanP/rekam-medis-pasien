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

func TestRoomTypeHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	expectedResponse := &dto.RoomTypeListResponse{
		Data: []dto.RoomTypeResponse{},
		Meta: dto.RoomTypePaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("List", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/room-types?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_DeletedList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	expectedResponse := &dto.RoomTypeDeletedListResponse{
		Data: []dto.DeletedRoomTypeResponse{},
		Meta: dto.RoomTypePaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("DeletedList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/room-types/deleted?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.DeletedList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_ActiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	expectedResponse := &dto.RoomTypeListResponse{
		Data: []dto.RoomTypeResponse{},
		Meta: dto.RoomTypePaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("ActiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/room-types/active?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ActiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_InactiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	expectedResponse := &dto.RoomTypeListResponse{
		Data: []dto.RoomTypeResponse{},
		Meta: dto.RoomTypePaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("InactiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/room-types/inactive?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.InactiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_FindByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	expectedResponse := mocks.NewTestRoomTypeResponseWithData(1, "ICU", "ICU-01", "Intensive Care Unit", true)

	mockService.On("FindByID", uint(1)).Return(expectedResponse, nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/room-types/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.FindByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	createReq := mocks.NewCreateRoomTypeRequest("ICU", "ICU-01", "Intensive Care Unit", true)
	body, _ := json.Marshal(createReq)

	expectedResponse := mocks.NewTestRoomTypeResponseWithData(1, "ICU", "ICU-01", "Intensive Care Unit", true)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateRoomTypeRequest")).Return(expectedResponse, nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/room-types", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRoomTypeHandler_Create_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	createReq := mocks.NewCreateRoomTypeRequest("ICU", "ICU-01", "Intensive Care Unit", true)
	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateRoomTypeRequest")).Return(nil, errors.New("name already exists"))

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/room-types", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRoomTypeHandler_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	updateReq := mocks.NewUpdateRoomTypeRequest("ICU Updated", "ICU-02", "Updated", true)
	body, _ := json.Marshal(updateReq)

	expectedResponse := mocks.NewTestRoomTypeResponseWithData(1, "ICU Updated", "ICU-02", "Updated", true)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateRoomTypeRequest")).Return(expectedResponse, nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/room-types/1", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_SoftDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	mockService.On("SoftDelete", uint(1)).Return(nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/room-types/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_Restore_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	mockService.On("Restore", uint(1)).Return(nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/room-types/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Restore(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_HardDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	mockService.On("HardDelete", uint(1)).Return(nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/room-types/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_Activate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	mockService.On("Activate", uint(1)).Return(nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/room-types/1/activate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Activate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomTypeHandler_Deactivate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockRoomTypeService)

	mockService.On("Deactivate", uint(1)).Return(nil)

	h := handler.NewRoomTypeHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/room-types/1/deactivate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Deactivate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
