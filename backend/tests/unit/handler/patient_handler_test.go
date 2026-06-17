package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============= Test Cases: GetMyPatientData =============

func TestGetMyPatientData_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	expectedPatient := mocks.NewTestPatientResponse()
	mockService.On("GetMyPatientData", uint(1)).Return(expectedPatient, nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patients/me", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.GetMyPatientData(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Patient data retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestGetMyPatientData_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)
	h := handler.NewPatientHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patients/me", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetMyPatientData(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User not authenticated", response["message"])
}

// ============= Test Cases: UpdateMyPatientData =============

func TestUpdateMyPatientData_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	updatedPatient := mocks.NewTestPatientResponse()
	updatedPhone := "08111111111"
	updatedPatient.Phone = updatedPhone

	mockService.On("UpdateMyPatientData", uint(1), mock.Anything).Return(updatedPatient, nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/patients/me", strings.NewReader(`{
		"phone": "08111111111"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.UpdateMyPatientData(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateMyPatientData_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	mockService.On("UpdateMyPatientData", uint(1), mock.Anything).Return(nil, errors.New("patient not found"))

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()

	req := httptest.NewRequest("PUT", "/patients/me", strings.NewReader(`{
		"phone": "08111111111"
	}`))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	h.UpdateMyPatientData(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: ListPatients =============

func TestListPatients_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	expectedResponse := &dto.PatientListResponse{
		Data: []dto.PatientResponse{*mocks.NewTestPatientResponse()},
		Meta: dto.PatientPaginationMeta{Page: 1, PageSize: 10, TotalItems: 1, TotalPages: 1},
	}

	mockService.On("ListPatients", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patients?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListPatients(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Patients retrieved successfully", response["message"])
}

func TestListPatients_InvalidPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)
	h := handler.NewPatientHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patients?page=invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListPatients(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============= Test Cases: GetPatientByID =============

func TestGetPatientByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	expectedPatient := mocks.NewTestPatientResponse()
	mockService.On("GetPatientByID", uint(1)).Return(expectedPatient, nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patients/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetPatientByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetPatientByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)
	h := handler.NewPatientHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patients/invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.GetPatientByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid patient ID", response["message"])
}

func TestGetPatientByID_OwnershipForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	userID := uint(2)
	patient := mocks.NewTestPatientResponse()
	patient.UserID = &userID
	mockService.On("GetPatientByID", uint(1)).Return(patient, nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patients/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("user_id", uint(1))
	c.Set("role", models.RolePatient)

	h.GetPatientByID(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: CreatePatient =============

func TestCreatePatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	expectedPatient := mocks.NewTestPatientResponse()
	mockService.On("CreatePatient", mock.Anything).Return(expectedPatient, nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/patients", strings.NewReader(`{
		"patient_code": "PAT001",
		"full_name": "John Doe",
		"date_of_birth": "1990-01-01",
		"gender": "male",
		"blood_type": "O"
	}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreatePatient(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreatePatient_DuplicateCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	mockService.On("CreatePatient", mock.Anything).Return(nil, errors.New("patient code already exists"))

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/patients", strings.NewReader(`{
		"patient_code": "PAT001",
		"full_name": "John Doe",
		"date_of_birth": "1990-01-01",
		"gender": "male",
		"blood_type": "O"
	}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreatePatient(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: UpdatePatient / Delete / Restore =============

func TestUpdatePatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)

	expectedPatient := mocks.NewTestPatientResponse()
	mockService.On("UpdatePatient", uint(1), mock.Anything).Return(expectedPatient, nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/patients/1", strings.NewReader(`{
		"full_name": "John Updated"
	}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UpdatePatient(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSoftDeletePatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)
	mockService.On("SoftDeletePatient", uint(1)).Return(nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/patients/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDeletePatient(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestRestorePatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)
	mockService.On("RestorePatient", uint(1)).Return(nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/patients/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.RestorePatient(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHardDeletePatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPatientService)
	mockService.On("HardDeletePatient", uint(1)).Return(nil)

	h := handler.NewPatientHandler(mockService)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/patients/1/permanent", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDeletePatient(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
