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

func TestTypeTestHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	expectedResponse := &dto.TypeTestListResponse{
		Data: []dto.TypeTestResponse{},
		Meta: dto.TypeTestPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("List", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-tests?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_ActiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	expectedResponse := &dto.TypeTestListResponse{
		Data: []dto.TypeTestResponse{},
		Meta: dto.TypeTestPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("ActiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-tests/active?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ActiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_InactiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	expectedResponse := &dto.TypeTestListResponse{
		Data: []dto.TypeTestResponse{},
		Meta: dto.TypeTestPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("InactiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-tests/inactive?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.InactiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_DeletedList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	expectedResponse := &dto.TypeTestDeletedListResponse{
		Data: []dto.DeletedTypeTestResponse{},
		Meta: dto.TypeTestPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("DeletedList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-tests/deleted?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.DeletedList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_FindByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	expectedResponse := mocks.NewTestTypeTestResponseWithData(1, "Test1", "T001", 1, "Desc1", 50000, true)

	mockService.On("FindByID", uint(1)).Return(expectedResponse, nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-tests/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.FindByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	createReq := mocks.NewCreateTypeTestRequest("Test1", "T001", 1, "Desc1", 50000, true)
	body, _ := json.Marshal(createReq)

	expectedResponse := mocks.NewTestTypeTestResponseWithData(1, "Test1", "T001", 1, "Desc1", 50000, true)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateTypeTestRequest")).Return(expectedResponse, nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/type-tests", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestTypeTestHandler_Create_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	createReq := mocks.NewCreateTypeTestRequest("Test1", "T001", 1, "Desc1", 50000, true)
	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateTypeTestRequest")).Return(nil, errors.New("code already exists"))

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/type-tests", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestTypeTestHandler_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	updateReq := mocks.NewUpdateTypeTestRequest("Test2", "T002", 2, "Desc2", 60000, true)
	body, _ := json.Marshal(updateReq)

	expectedResponse := mocks.NewTestTypeTestResponseWithData(1, "Test2", "T002", 2, "Desc2", 60000, true)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateTypeTestRequest")).Return(expectedResponse, nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/type-tests/1", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_SoftDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	mockService.On("SoftDelete", uint(1)).Return(nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/type-tests/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_Restore_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	mockService.On("Restore", uint(1)).Return(nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/type-tests/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Restore(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_HardDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	mockService.On("HardDelete", uint(1)).Return(nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/type-tests/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_Activate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	mockService.On("Activate", uint(1)).Return(nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/type-tests/1/activate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Activate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestHandler_Deactivate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestService)

	mockService.On("Deactivate", uint(1)).Return(nil)

	h := handler.NewTypeTestHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/type-tests/1/deactivate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Deactivate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
