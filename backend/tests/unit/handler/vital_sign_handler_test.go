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

func setupVitalSignHandlerTest() (*gin.Engine, *mocks.MockVitalSignService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockVitalSignService)

	h := handler.NewVitalSignHandler(mockService)

	r := gin.Default()

	r.GET("/vital-signs", h.List)
	r.GET("/vital-signs/deleted", h.DeletedList)
	r.GET("/vital-signs/:id", h.FindByID)
	r.POST("/vital-signs", h.Create)
	r.PUT("/vital-signs/:id", h.Update)
	r.DELETE("/vital-signs/:id", h.SoftDelete)
	r.PATCH("/vital-signs/:id/restore", h.Restore)
	r.DELETE("/vital-signs/:id/hard", h.HardDelete)

	return r, mockService
}

func TestVitalSignHandler_List_Success(t *testing.T) {
	r, mockService := setupVitalSignHandlerTest()

	expectedRes := &dto.VitalSignListResponse{
		Data: []dto.VitalSignResponse{
			*mocks.NewTestVitalSignResponse(mocks.NewTestVitalSignWithData(1, 1, false)),
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.VitalSignPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/vital-signs?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestVitalSignHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupVitalSignHandlerTest()

	expectedRes := mocks.NewTestVitalSignResponse(mocks.NewTestVitalSignWithData(1, 1, false))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/vital-signs/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestVitalSignHandler_Create_Success(t *testing.T) {
	r, mockService := setupVitalSignHandlerTest()
	createReq := mocks.NewCreateVitalSignRequest(1)
	expectedRes := mocks.NewTestVitalSignResponse(mocks.NewTestVitalSignWithData(1, 1, false))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateVitalSignRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/vital-signs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestVitalSignHandler_Update_Success(t *testing.T) {
	r, mockService := setupVitalSignHandlerTest()
	updateReq := mocks.NewUpdateVitalSignRequest()
	expectedRes := mocks.NewTestVitalSignResponse(mocks.NewTestVitalSignWithData(1, 1, false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateVitalSignRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/vital-signs/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestVitalSignHandler_SoftDelete_Success(t *testing.T) {
	r, mockService := setupVitalSignHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/vital-signs/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestVitalSignHandler_Restore_Success(t *testing.T) {
	r, mockService := setupVitalSignHandlerTest()

	mockService.On("Restore", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/vital-signs/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestVitalSignHandler_HardDelete_Success(t *testing.T) {
	r, mockService := setupVitalSignHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/vital-signs/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
