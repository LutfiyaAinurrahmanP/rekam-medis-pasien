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

func setupLabTestHandlerTest() (*gin.Engine, *mocks.MockLabTestService, *mocks.MockDoctorRepository) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockLabTestService)
	mockDoctorRepo := new(mocks.MockDoctorRepository)

	h := handler.NewLabTestHandler(mockService, mockDoctorRepo)

	r := gin.Default()

	// Add test middleware for role
	r.Use(func(c *gin.Context) {
		c.Set("role", "doctor")
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.GET("/lab-tests", h.List)
	r.GET("/lab-tests/deleted", h.DeletedList)
	r.GET("/lab-tests/:id", h.FindByID)
	r.GET("/lab-tests/medical-record/:record_id", h.FindByMedicalRecordID)
	r.POST("/lab-tests", h.Create)
	r.PUT("/lab-tests/:id", h.Update)
	r.PATCH("/lab-tests/:id/collect-sample", h.CollectSample)
	r.PATCH("/lab-tests/:id/start", h.Start)
	r.PATCH("/lab-tests/:id/complete", h.Complete)
	r.PATCH("/lab-tests/:id/cancel", h.Cancel)
	r.DELETE("/lab-tests/:id", h.SoftDelete)
	r.PATCH("/lab-tests/:id/restore", h.Restore)
	r.DELETE("/lab-tests/:id/hard", h.HardDelete)

	return r, mockService, mockDoctorRepo
}

func TestLabTestHandler_List_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()

	expectedRes := &dto.LabTestListResponse{
		Data: []dto.LabTestResponse{
			*mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)),
		},
		Meta: dto.LabTestPaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.LabTestPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/lab-tests?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_FindByID_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()

	expectedRes := mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/lab-tests/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_FindByMedicalRecordID_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()

	expectedRes := &dto.LabTestListResponse{
		Data: []dto.LabTestResponse{
			*mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)),
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.LabTestPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/lab-tests/medical-record/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_Create_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()
	createReq := mocks.NewCreateLabTestRequest(1, 1, 1)
	expectedRes := mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateLabTestRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/lab-tests", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_Update_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()
	updateReq := mocks.NewUpdateLabTestRequest("Notes", "10-20")
	expectedRes := mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateLabTestRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/lab-tests/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_CollectSample_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()
	expectedRes := mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "sample_collected", false))

	mockService.On("CollectSample", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/lab-tests/1/collect-sample", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_Start_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()
	expectedRes := mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "in_progress", false))

	mockService.On("Start", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/lab-tests/1/start", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_Complete_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()
	completeReq := &dto.CompleteLabTestRequest{}
	expectedRes := mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "completed", false))

	body, _ := json.Marshal(completeReq)

	mockService.On("Complete", uint(1), mock.AnythingOfType("*dto.CompleteLabTestRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/lab-tests/1/complete", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_Cancel_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()
	expectedRes := mocks.NewTestLabTestResponse(mocks.NewTestLabTestWithData(1, 1, 1, 1, "cancelled", false))

	mockService.On("Cancel", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/lab-tests/1/cancel", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_SoftDelete_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/lab-tests/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_Restore_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()

	mockService.On("Restore", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/lab-tests/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLabTestHandler_HardDelete_Success(t *testing.T) {
	r, mockService, _ := setupLabTestHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/lab-tests/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
