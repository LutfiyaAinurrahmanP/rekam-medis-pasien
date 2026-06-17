package handler

import (
	"bytes"
	"encoding/json"
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

func setupFamilyHistoryHandlerTest() (*gin.Engine, *mocks.MockFamilyHistoryService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockFamilyHistoryService)

	h := handler.NewFamilyHistoryHandler(mockService)

	r := gin.Default()

	r.GET("/medical-history/family-histories", h.List)
	r.GET("/medical-history/family-histories/:id", h.FindByID)
	r.POST("/medical-history/family-histories", h.Create)
	r.PUT("/medical-history/family-histories/:id", h.Update)
	r.DELETE("/medical-history/family-histories/:id", h.Delete)

	return r, mockService
}

func TestFamilyHistoryHandler_List_Success(t *testing.T) {
	r, mockService := setupFamilyHistoryHandlerTest()

	expectedRes := &dto.FamilyHistoryListResponse{
		Data: []dto.FamilyHistoryResponse{
			*mocks.NewTestFamilyHistoryResponse(mocks.NewTestFamilyHistoryWithData(1, 1)),
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.FamilyHistoryPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/family-histories?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestFamilyHistoryHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupFamilyHistoryHandlerTest()

	expectedRes := mocks.NewTestFamilyHistoryResponse(mocks.NewTestFamilyHistoryWithData(1, 1))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/family-histories/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestFamilyHistoryHandler_Create_Success(t *testing.T) {
	r, mockService := setupFamilyHistoryHandlerTest()
	createReq := mocks.NewCreateFamilyHistoryRequest(1)
	expectedRes := mocks.NewTestFamilyHistoryResponse(mocks.NewTestFamilyHistoryWithData(1, 1))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateFamilyHistoryRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/medical-history/family-histories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestFamilyHistoryHandler_Update_Success(t *testing.T) {
	r, mockService := setupFamilyHistoryHandlerTest()
	updateReq := mocks.NewUpdateFamilyHistoryRequest()
	expectedRes := mocks.NewTestFamilyHistoryResponse(mocks.NewTestFamilyHistoryWithData(1, 1))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateFamilyHistoryRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/medical-history/family-histories/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestFamilyHistoryHandler_Delete_Success(t *testing.T) {
	r, mockService := setupFamilyHistoryHandlerTest()

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medical-history/family-histories/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
