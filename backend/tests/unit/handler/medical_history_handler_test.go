package handler

import (
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

func setupMedicalHistoryHandlerTest() (*gin.Engine, *mocks.MockMedicalHistoryService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockMedicalHistoryService)

	h := handler.NewMedicalHistoryHandler(mockService)

	r := gin.Default()

	r.GET("/medical-history", h.List)
	r.GET("/medical-history/:id", h.FindByID)
	r.GET("/medical-history/patient/:pid", h.FindByPatientID)

	return r, mockService
}

func TestMedicalHistoryHandler_List_Success(t *testing.T) {
	r, mockService := setupMedicalHistoryHandlerTest()

	expectedRes := &dto.MedicalHistoryListResponse{
		Data: []dto.MedicalHistoryOverviewResponse{
			{ID: 1, PatientID: 1, PatientName: "Budi Santoso"},
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.MedicalHistoryPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalHistoryHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupMedicalHistoryHandlerTest()

	expectedRes := &dto.MedicalHistoryDetailResponse{
		ID:        1,
		PatientID: 1,
	}

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalHistoryHandler_FindByPatientID_Success(t *testing.T) {
	r, mockService := setupMedicalHistoryHandlerTest()

	expectedRes := &dto.MedicalHistoryDetailResponse{
		ID:        1,
		PatientID: 1,
	}

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/patient/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
