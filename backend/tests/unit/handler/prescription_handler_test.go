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

func setupPrescriptionHandlerTest() (*gin.Engine, *mocks.MockPrescriptionService, *mocks.MockDoctorRepository, *mocks.MockPatientRepository) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPrescriptionService)
	mockDoctorRepo := new(mocks.MockDoctorRepository)
	mockPatientRepo := new(mocks.MockPatientRepository)

	h := handler.NewPrescriptionHandler(mockService, mockDoctorRepo, mockPatientRepo)

	r := gin.Default()

	// Add test middleware for role
	r.Use(func(c *gin.Context) {
		c.Set("role", "doctor")
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.GET("/prescriptions", h.List)
	r.GET("/prescriptions/deleted", h.DeletedList)
	r.GET("/prescriptions/:id", h.FindByID)
	r.GET("/prescriptions/medical-record/:recordID", h.PrescriptionsByMedicalRecordID)
	r.POST("/prescriptions", h.Create)
	r.PUT("/prescriptions/:id", h.Update)
	r.PATCH("/prescriptions/:id/dispense", h.Dispense)
	r.PATCH("/prescriptions/:id/cancel", h.Cancel)
	r.DELETE("/prescriptions/:id", h.SoftDelete)
	r.PATCH("/prescriptions/:id/restore", h.Restore)
	r.DELETE("/prescriptions/:id/hard", h.HardDelete)

	r.GET("/prescriptions/:id/items", h.ListItems)
	r.GET("/prescriptions/:id/items/:itemID", h.FindItemByID)
	r.POST("/prescriptions/:id/items", h.CreateItem)
	r.PUT("/prescriptions/:id/items/:itemID", h.UpdateItem)
	r.DELETE("/prescriptions/:id/items/:itemID", h.DeleteItem)

	return r, mockService, mockDoctorRepo, mockPatientRepo
}

func TestPrescriptionHandler_List_Success(t *testing.T) {
	r, mockService, mockDoctorRepo, _ := setupPrescriptionHandlerTest()

	expectedRes := &dto.PrescriptionListResponse{
		Data: []dto.PrescriptionResponse{
			*mocks.NewTestPrescriptionResponse(mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)),
		},
	}

	uid := uint(1)
	mockDoctorRepo.On("FindByUserID", uint(1)).Return(&models.Doctor{ID: 1, UserID: &uid}, nil)
	mockService.On("List", mock.AnythingOfType("*dto.PrescriptionPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/prescriptions?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	mockDoctorRepo.AssertExpectations(t)
}

func TestPrescriptionHandler_FindByID_Success(t *testing.T) {
	r, mockService, mockDoctorRepo, _ := setupPrescriptionHandlerTest()

	expectedRes := mocks.NewTestPrescriptionResponse(mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false))

	uid := uint(1)
	mockDoctorRepo.On("FindByUserID", uint(1)).Return(&models.Doctor{ID: 1, UserID: &uid}, nil)
	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/prescriptions/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestPrescriptionHandler_Create_Success(t *testing.T) {
	r, mockService, mockDoctorRepo, _ := setupPrescriptionHandlerTest()
	createReq := mocks.NewCreatePrescriptionRequest(1, 1)
	expectedRes := mocks.NewTestPrescriptionResponse(mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false))

	body, _ := json.Marshal(createReq)

	uid := uint(1)
	mockDoctorRepo.On("FindByUserID", uint(1)).Return(&models.Doctor{ID: 1, UserID: &uid}, nil)
	mockService.On("Create", mock.AnythingOfType("*dto.CreatePrescriptionRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/prescriptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestPrescriptionHandler_Update_Success(t *testing.T) {
	r, mockService, _, _ := setupPrescriptionHandlerTest()
	updateReq := mocks.NewUpdatePrescriptionRequest("Notes")
	expectedRes := mocks.NewTestPrescriptionResponse(mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdatePrescriptionRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/prescriptions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestPrescriptionHandler_Dispense_Success(t *testing.T) {
	r, mockService, _, _ := setupPrescriptionHandlerTest()
	expectedRes := mocks.NewTestPrescriptionResponse(mocks.NewTestPrescriptionWithData(1, 1, 1, "dispensed", false))

	mockService.On("Dispense", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPatch, "/prescriptions/1/dispense", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestPrescriptionHandler_CreateItem_Success(t *testing.T) {
	r, mockService, _, _ := setupPrescriptionHandlerTest()
	createReq := mocks.NewCreatePrescriptionItemRequest(1)
	expectedRes := mocks.NewTestPrescriptionItemResponse(mocks.NewTestPrescriptionItemWithData(1, 1, 1))

	body, _ := json.Marshal(createReq)

	mockService.On("CreateItem", uint(1), mock.AnythingOfType("*dto.CreatePrescriptionItemRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/prescriptions/1/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestPrescriptionHandler_UpdateItem_Success(t *testing.T) {
	r, mockService, _, _ := setupPrescriptionHandlerTest()
	updateReq := mocks.NewUpdatePrescriptionItemRequest("2 tablets")
	expectedRes := mocks.NewTestPrescriptionItemResponse(mocks.NewTestPrescriptionItemWithData(1, 1, 1))

	body, _ := json.Marshal(updateReq)

	mockService.On("UpdateItem", uint(1), uint(1), mock.AnythingOfType("*dto.UpdatePrescriptionItemRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/prescriptions/1/items/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestPrescriptionHandler_DeleteItem_Success(t *testing.T) {
	r, mockService, _, _ := setupPrescriptionHandlerTest()

	mockService.On("DeleteItem", uint(1), uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/prescriptions/1/items/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
