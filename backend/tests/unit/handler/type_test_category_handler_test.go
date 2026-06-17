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

func TestTypeTestCategoryHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	expectedResponse := &dto.TypeTestCategoryListResponse{
		Data: []dto.TypeTestCategoryResponse{},
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("List", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-test-categories?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_ActiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	expectedResponse := &dto.TypeTestCategoryListResponse{
		Data: []dto.TypeTestCategoryResponse{},
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("ActiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-test-categories/active?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ActiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_InactiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	expectedResponse := &dto.TypeTestCategoryListResponse{
		Data: []dto.TypeTestCategoryResponse{},
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("InactiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-test-categories/inactive?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.InactiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_DeletedList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	expectedResponse := &dto.TypeTestCategoryDeletedListResponse{
		Data: []dto.DeletedTypeTestCategoryResponse{},
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("DeletedList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-test-categories/deleted?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.DeletedList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_FindByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	expectedResponse := mocks.NewTestTypeTestCategoryResponseWithData(1, "Cat1", "C001", "Desc1", true)

	mockService.On("FindByID", uint(1)).Return(expectedResponse, nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/type-test-categories/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.FindByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	createReq := mocks.NewCreateTypeTestCategoryRequest("Cat1", "C001", "Desc1", true)
	body, _ := json.Marshal(createReq)

	expectedResponse := mocks.NewTestTypeTestCategoryResponseWithData(1, "Cat1", "C001", "Desc1", true)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateTypeTestCategoryRequest")).Return(expectedResponse, nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/type-test-categories", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestTypeTestCategoryHandler_Create_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	createReq := mocks.NewCreateTypeTestCategoryRequest("Cat1", "C001", "Desc1", true)
	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateTypeTestCategoryRequest")).Return(nil, errors.New("name already exists"))

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/type-test-categories", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestTypeTestCategoryHandler_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	updateReq := mocks.NewUpdateTypeTestCategoryRequest("Cat2", "C002", "Desc2", true)
	body, _ := json.Marshal(updateReq)

	expectedResponse := mocks.NewTestTypeTestCategoryResponseWithData(1, "Cat2", "C002", "Desc2", true)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateTypeTestCategoryRequest")).Return(expectedResponse, nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/type-test-categories/1", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_SoftDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	mockService.On("SoftDelete", uint(1)).Return(nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/type-test-categories/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_Restore_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	mockService.On("Restore", uint(1)).Return(nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/type-test-categories/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Restore(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_HardDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	mockService.On("HardDelete", uint(1)).Return(nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/type-test-categories/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_Activate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	mockService.On("Activate", uint(1)).Return(nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/type-test-categories/1/activate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Activate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTypeTestCategoryHandler_Deactivate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockTypeTestCategoryService)

	mockService.On("Deactivate", uint(1)).Return(nil)

	h := handler.NewTypeTestCategoryHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/type-test-categories/1/deactivate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Deactivate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
