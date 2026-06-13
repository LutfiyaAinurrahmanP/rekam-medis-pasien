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

func setupSurgicalHistoryHandlerTest() (*gin.Engine, *mocks.MockSurgicalHistoryService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockSurgicalHistoryService)

	h := handler.NewSurgicalHistoryHandler(mockService)

	r := gin.Default()

	r.GET("/medical-history/surgeries", h.List)
	r.GET("/medical-history/surgeries/:id", h.FindByID)
	r.POST("/medical-history/surgeries", h.Create)
	r.PUT("/medical-history/surgeries/:id", h.Update)
	r.DELETE("/medical-history/surgeries/:id", h.Delete)

	return r, mockService
}

func TestSurgicalHistoryHandler_List_Success(t *testing.T) {
	r, mockService := setupSurgicalHistoryHandlerTest()

	expectedRes := &dto.SurgicalHistoryListResponse{
		Data: []dto.SurgicalHistoryResponse{
			*mocks.NewTestSurgicalHistoryResponse(mocks.NewTestSurgicalHistoryWithData(1, 1)),
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.SurgicalHistoryPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/surgeries?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSurgicalHistoryHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupSurgicalHistoryHandlerTest()

	expectedRes := mocks.NewTestSurgicalHistoryResponse(mocks.NewTestSurgicalHistoryWithData(1, 1))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/surgeries/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSurgicalHistoryHandler_Create_Success(t *testing.T) {
	r, mockService := setupSurgicalHistoryHandlerTest()
	createReq := mocks.NewCreateSurgicalHistoryRequest(1)
	expectedRes := mocks.NewTestSurgicalHistoryResponse(mocks.NewTestSurgicalHistoryWithData(1, 1))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateSurgicalHistoryRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/medical-history/surgeries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestSurgicalHistoryHandler_Update_Success(t *testing.T) {
	r, mockService := setupSurgicalHistoryHandlerTest()
	updateReq := mocks.NewUpdateSurgicalHistoryRequest()
	expectedRes := mocks.NewTestSurgicalHistoryResponse(mocks.NewTestSurgicalHistoryWithData(1, 1))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateSurgicalHistoryRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/medical-history/surgeries/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSurgicalHistoryHandler_Delete_Success(t *testing.T) {
	r, mockService := setupSurgicalHistoryHandlerTest()

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medical-history/surgeries/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
