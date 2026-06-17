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

func setupHospitalizationHandlerTest() (*gin.Engine, *mocks.MockHospitalizationService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockHospitalizationService)

	h := handler.NewHospitalizationHandler(mockService)

	r := gin.Default()

	r.GET("/hospitalizations", h.List)
	r.GET("/hospitalizations/deleted", h.DeletedList)
	r.GET("/hospitalizations/active", h.ActiveList)
	r.GET("/hospitalizations/inactive", h.InactiveList)
	r.GET("/hospitalizations/:id", h.FindByID)
	r.POST("/hospitalizations", h.Create)
	r.PUT("/hospitalizations/:id", h.Update)
	r.PATCH("/hospitalizations/:id/discharge", h.Discharge)
	r.PATCH("/hospitalizations/:id/transfer", h.Transfer)
	r.PATCH("/hospitalizations/:id/activate", h.Activate)
	r.PATCH("/hospitalizations/:id/deactivate", h.Deactivate)
	r.DELETE("/hospitalizations/:id", h.SoftDelete)
	r.PATCH("/hospitalizations/:id/restore", h.Restore)
	r.DELETE("/hospitalizations/:id/hard", h.HardDelete)

	return r, mockService
}

func TestHospitalizationHandler_List_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()

	expectedRes := &dto.HospitalizationListResponse{
		Data: []dto.HospitalizationResponse{
			*mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)),
		},
		Meta: dto.HospitalizationPaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.HospitalizationPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/hospitalizations?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()

	expectedRes := mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/hospitalizations/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_Create_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()
	createReq := mocks.NewCreateHospitalizationRequest(1, 1, 1, "2023-12-01T10:00:00Z", "Reason")
	expectedRes := mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateHospitalizationRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/hospitalizations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_Update_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()
	updateReq := mocks.NewUpdateHospitalizationRequest(2, 2, "New Reason")
	expectedRes := mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 2, 2, "2023-12-01", "10:00:00", "admitted", false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateHospitalizationRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/hospitalizations/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_Discharge_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()
	dischargeReq := &dto.DischargeHospitalizationRequest{DischargeSummary: "Improved"}
	expectedRes := mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "discharged", false))

	body, _ := json.Marshal(dischargeReq)

	mockService.On("Discharge", uint(1), mock.AnythingOfType("*dto.DischargeHospitalizationRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/hospitalizations/1/discharge", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_Transfer_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()
	transferReq := &dto.TransferHospitalizationRequest{Notes: "To ICU"}
	expectedRes := mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "transferred", false))

	body, _ := json.Marshal(transferReq)

	mockService.On("Transfer", uint(1), mock.AnythingOfType("*dto.TransferHospitalizationRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/hospitalizations/1/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_Activate_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()
	expectedRes := mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false))

	mockService.On("Activate", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/hospitalizations/1/activate", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_Deactivate_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()
	expectedRes := mocks.NewTestHospitalizationResponse(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "discharged", false))

	mockService.On("Deactivate", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/hospitalizations/1/deactivate", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_SoftDelete_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/hospitalizations/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_Restore_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()

	mockService.On("Restore", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/hospitalizations/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHospitalizationHandler_HardDelete_Success(t *testing.T) {
	r, mockService := setupHospitalizationHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/hospitalizations/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
