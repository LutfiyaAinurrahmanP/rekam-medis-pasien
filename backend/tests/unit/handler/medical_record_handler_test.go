package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMedicalRecordHandlerTest() (*gin.Engine, *mocks.MockMedicalRecordService, *mocks.MockDoctorRepository, *mocks.MockPatientRepository) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockMedicalRecordService)
	mockDoctorRepo := new(mocks.MockDoctorRepository)
	mockPatientRepo := new(mocks.MockPatientRepository)

	h := handler.NewMedicalRecordHandler(mockService, mockDoctorRepo, mockPatientRepo)

	r := gin.Default()

	// Inject role patient manually for testing ownership
	r.Use(func(c *gin.Context) {
		c.Set("role", "patient")
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.GET("/medical-records/me", h.MyMedicalRecords)
	r.GET("/medical-records/patient/:patientID", h.MedicalRecordsByPatientID)
	r.GET("/medical-records", h.List)
	r.GET("/medical-records/deleted", h.DeletedList)
	r.GET("/medical-records/:id", h.FindByID)
	r.POST("/medical-records", h.Create)
	r.PUT("/medical-records/:id", h.Update)
	r.PATCH("/medical-records/:id/finalize", h.Finalize)
	r.DELETE("/medical-records/:id", h.SoftDelete)
	r.PATCH("/medical-records/:id/restore", h.Restore)
	r.DELETE("/medical-records/:id/hard", h.HardDelete)

	return r, mockService, mockDoctorRepo, mockPatientRepo
}

func TestMedicalRecordHandler_List_Success(t *testing.T) {
	r, mockService, _, _ := setupMedicalRecordHandlerTest()

	expectedRes := &dto.MedicalRecordListResponse{
		Data: []dto.MedicalRecordResponse{
			*mocks.NewTestMedicalRecordResponse(mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false)),
		},
		Meta: dto.MedicalRecordPaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.MedicalRecordPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-records?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalRecordHandler_FindByID_Success(t *testing.T) {
	r, mockService, _, mockPatientRepo := setupMedicalRecordHandlerTest()

	expectedRes := mocks.NewTestMedicalRecordResponse(mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false))

	// Mock for ownership check
	uid := uint(1)
	mockPatientRepo.On("FindByUserID", uint(1)).Return(&models.Patient{ID: 1, UserID: &uid}, nil)
	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-records/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	mockPatientRepo.AssertExpectations(t)
}

func TestMedicalRecordHandler_Create_Success(t *testing.T) {
	r, mockService, _, _ := setupMedicalRecordHandlerTest()
	createReq := mocks.NewCreateMedicalRecordRequest(1, "2023-12-01", "Complaint", "Diagnosis", "Plan")
	expectedRes := mocks.NewTestMedicalRecordResponse(mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateMedicalRecordRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/medical-records", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalRecordHandler_Update_Success(t *testing.T) {
	r, mockService, _, _ := setupMedicalRecordHandlerTest()
	updateReq := mocks.NewUpdateMedicalRecordRequest("New Complaint", "New Diagnosis", "New Plan")
	expectedRes := mocks.NewTestMedicalRecordResponse(mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "New Complaint", "New Diagnosis", "New Plan", "draft", false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateMedicalRecordRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/medical-records/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalRecordHandler_Finalize_Success(t *testing.T) {
	r, mockService, _, _ := setupMedicalRecordHandlerTest()

	mockService.On("Finalize", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/medical-records/1/finalize", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalRecordHandler_SoftDelete_Success(t *testing.T) {
	r, mockService, _, _ := setupMedicalRecordHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medical-records/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalRecordHandler_Restore_Success(t *testing.T) {
	r, mockService, _, _ := setupMedicalRecordHandlerTest()

	mockService.On("Restore", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/medical-records/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalRecordHandler_HardDelete_Success(t *testing.T) {
	r, mockService, _, _ := setupMedicalRecordHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medical-records/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
